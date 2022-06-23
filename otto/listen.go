package otto

import (
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

// ListenOptions describes options for listening
type ListenOptions struct {
	Address           string
	AllowFrom         []net.IPNet
	Identity          ssh.Signer
	TrustedPublicKeys []string
}

// Listener describes an active listening Otto server
type Listener struct {
	options ListenOptions
	handle  func(conn *Connection)
	l       net.Listener
}

// SetupListener will prepare a listening socket for incoming connections. No connections are accepted until you call
// Accept().
func SetupListener(options ListenOptions, handle func(conn *Connection)) (*Listener, error) {
	for _, trustedKey := range options.TrustedPublicKeys {
		if trustedKey == base64.StdEncoding.EncodeToString(options.Identity.PublicKey().Marshal()) {
			return nil, fmt.Errorf("server and client identity cannot be the same")
		}
	}

	l, err := net.Listen("tcp", options.Address)
	if err != nil {
		return nil, err
	}
	log.Info("Otto client listening on %s", options.Address)
	return &Listener{
		options: options,
		handle:  handle,
		l:       l,
	}, nil
}

// Port get the port the listener is listening on
func (l *Listener) Port() uint16 {
	p := strings.Split(l.l.Addr().String(), ":")
	port, err := strconv.ParseUint(p[len(p)-1], 10, 16)
	if err != nil {
		panic("invalid port")
	}
	return uint16(port)
}

// Accept will accpet incoming connections. Blocking.
func (l *Listener) Accept() error {
	for {
		c, err := l.l.Accept()
		if err != nil {
			log.PDebug("[LISTEN] Error accepting connection", map[string]interface{}{
				"error": err.Error(),
			})
			return err
		}
		log.PDebug("[LISTEN] Incoming connection", map[string]interface{}{
			"remote_addr": c.RemoteAddr().String(),
		})
		go l.accept(c)
	}
}

// Close will stop the listener.
func (l *Listener) Close() {
	l.l.Close()
}

func (l *Listener) accept(c net.Conn) {
	connId := fdFromConn(c)

	if len(l.options.AllowFrom) > 0 {
		allow := false
		for _, allowNet := range l.options.AllowFrom {
			if allowNet.Contains(c.RemoteAddr().(*net.TCPAddr).IP) {
				log.PDebug("[LISTEN] Connection allowed by rule", map[string]interface{}{
					"remote_addr":     c.RemoteAddr().String(),
					"allowed_network": allowNet.String(),
				})
				allow = true
				break
			}
		}
		if !allow {
			log.PWarn("[LISTEN] Rejecting connection from server outside of allowed network", map[string]interface{}{
				"remote_addr":  c.RemoteAddr().String(),
				"allowed_addr": l.options.AllowFrom,
			})
			c.Close()
			return
		}
	}

	localIdentity := l.options.Identity.PublicKey().Marshal()
	var remoteIdentity []byte

	sshConfig := &ssh.ServerConfig{
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			log.PDebug("[LISTEN] Handshake", map[string]interface{}{
				"public_key": base64.StdEncoding.EncodeToString(pubKey.Marshal()),
			})
			remoteIdentity = pubKey.Marshal()

			for _, trustedKey := range l.options.TrustedPublicKeys {
				if trustedKey == base64.StdEncoding.EncodeToString(pubKey.Marshal()) {
					log.Debug("[LISTEN] Recognized public key")
					return &ssh.Permissions{
						Extensions: map[string]string{
							"pubkey-fp": ssh.FingerprintSHA256(pubKey),
						},
					}, nil
				}
			}
			log.PWarn("[LISTEN] Rejecting connection from untrusted public key", map[string]interface{}{
				"public_key": base64.StdEncoding.EncodeToString(pubKey.Marshal()),
			})
			return nil, fmt.Errorf("unknown public key %x", pubKey.Marshal())
		},
		ServerVersion: fmt.Sprintf("SSH-2.0-OTTO-%d", ProtocolVersion),
	}
	sshConfig.AddHostKey(l.options.Identity)

	sc, chans, reqs, err := ssh.NewServerConn(c, sshConfig)
	if err != nil {
		if err != io.EOF {
			log.PError("[LISTEN] SSH handshake error", map[string]interface{}{
				"remote_addr": c.RemoteAddr().String(),
				"error":       err.Error(),
			})
		}
		c.Close()
		return
	}

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		log.Debug("[LISTEN] ssh channel opened")
		if newChannel.ChannelType() != sshChannelName {
			log.PError("[LISTEN] Unknown SSH channel", map[string]interface{}{
				"channel_type": newChannel.ChannelType(),
				"remote_addr":  c.RemoteAddr().String(),
			})
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			return
		}
		channel, _, err := newChannel.Accept()
		if err != nil {
			log.PError("[LISTEN] SSH channel error", map[string]interface{}{
				"remote_addr": c.RemoteAddr().String(),
				"error":       err.Error(),
			})
			return
		}
		log.PDebug("[LISTEN] SSH handshake success", map[string]interface{}{
			"id":          connId,
			"remote_addr": c.RemoteAddr().String(),
		})
		l.handle(&Connection{
			id:             connId,
			w:              channel,
			remoteAddr:     c.RemoteAddr(),
			localAddr:      c.LocalAddr(),
			localIdentity:  localIdentity,
			remoteIdentity: remoteIdentity,
		})
		c.Close()
		channel.Close()
	}
	sc.Close()
}

func fdFromConn(c net.Conn) int {
	fdVal := reflect.Indirect(reflect.ValueOf(c)).FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	fdInt := int(pfdVal.FieldByName("Sysfd").Int())
	return fdInt
}
