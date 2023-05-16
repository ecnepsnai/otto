import * as React from 'react';
import { Host, HostType } from '../../types/Host';
import { Icon } from '../../components/Icon';
import { Badge } from '../../components/Badge';
import { Style } from '../../components/Style';
import { Dropdown, Menu } from '../../components/Menu';
import { Clipboard } from '../../services/Clipboard';
import { Notification } from '../../components/Notification';
import { GlobalModalFrame, Modal, ModalForm } from '../../components/Modal';
import { Pre } from '../../components/Pre';
import { Input } from '../../components/input/Input';
import { Permissions, UserAction } from '../../services/Permissions';

interface HostTrustProps {
    host: HostType;
    badgeOnly?: boolean;
    outline?: boolean;
    onReload?: () => (void);
}
export const HostTrust: React.FC<HostTrustProps> = (props: HostTrustProps) => {
    const addTrust = () => {
        let identity = '';
        const updateIdentity = (id: string) => {
            identity = id;
        };
        const saveIdentity = () => {
            return Host.UpdateTrust(props.host.ID, 'permit', identity).then(() => {
                Notification.success('Trust updated');
                props.onReload();
            });
        };

        GlobalModalFrame.showModal(
            <ModalForm title={'Add Trusted Identity'} onSubmit={saveIdentity}>
            <Input.Text
                label="Identity"
                type="text"
                defaultValue=""
                onChange={updateIdentity}
                required />
            </ModalForm>
        );
    };
    const removeTrust = () => {
        Modal.confirm('Remove Trust', 'Are you sure you want to remove this identity? The Otto server will be unable to communicate to this host until trust is reestablished.').then(confirmed => {
            if (!confirmed) {
                return;
            }
            Host.UpdateTrust(props.host.ID, 'deny').then(() => {
                Notification.success('Trust updated');
                props.onReload();
            });
        });
    };
    const trustPending = () => {
        const body = (<React.Fragment>
            <p>The following identity was detected on this host. Do you want to trust it?</p>
            <Pre>{props.host.Trust.UntrustedIdentity}</Pre>
        </React.Fragment>);
        Modal.confirm('Trust Pending Identity', body).then(confirmed => {
            if (!confirmed) {
                return;
            }
            Host.UpdateTrust(props.host.ID, 'permit').then(() => {
                Notification.success('Trust updated');
                props.onReload();
            });
        });
        
    };
    const copyServerIdentity = () => {
        Host.ServerID(props.host.ID).then(serverID => {
            Clipboard.setText(serverID).then(() => {
                Notification.success('Server ID Copied');
            });
        });
    };
    const rotateIdentity = () => {
        Host.RotateID(props.host.ID).then(() => {
            Notification.success('Identity Rotated');
        });
    };

    let badge = (<Badge pill outline={props.outline} color={Style.Palette.Secondary}><Icon.Label icon={<Icon.QuestionCircle />} label="Unknown" /></Badge>);

    if (props.host.Trust.TrustedIdentity) {
        badge = (<Badge pill outline={props.outline} color={Style.Palette.Success}><Icon.Label icon={<Icon.CheckCircle />} label="Established" /></Badge>);
    }
    if (props.host.Trust.UntrustedIdentity) {
        badge = (<Badge pill outline={props.outline} color={Style.Palette.Warning}><Icon.Label icon={<Icon.ExclamationTriangle />} label="Pending" /></Badge>);
    }

    if (props.badgeOnly) {
        return badge;
    }

    const setIdentityMenu = (<Menu.Item icon={<Icon.Plus />} label="Add Trusted Identity" onClick={addTrust} disabled={!Permissions.UserCan(UserAction.ModifyHosts)} />);
    const untrustIdentityMenu = props.host.Trust.TrustedIdentity ? (<Menu.Item icon={<Icon.Unlock />} label="Remove Trusted Identity" onClick={removeTrust} disabled={!Permissions.UserCan(UserAction.ModifyHosts)} />) : null;
    const trustPendingMenu = props.host.Trust.UntrustedIdentity ? (<Menu.Item icon={<Icon.Lock />} label="Confirm Pending Identity" onClick={trustPending} disabled={!Permissions.UserCan(UserAction.ModifyHosts)} />) : null;
    const copyServerIdentityMenu = (<Menu.Item icon={<Icon.Clipboard />} label="Copy Server Identity" onClick={copyServerIdentity} />);
    const rotateIdentityMenu = props.host.Trust.TrustedIdentity ? (<Menu.Item icon={<Icon.Random />} label="Rotate Identity" onClick={rotateIdentity} disabled={!Permissions.UserCan(UserAction.ModifyHosts)} />) : null;

    return (
        <span className="badges">
            {badge}
            <Dropdown label={<Icon.Bars />}>
                {setIdentityMenu}
                {untrustIdentityMenu}
                {trustPendingMenu}
                <Menu.Divider />
                {rotateIdentityMenu}
                {copyServerIdentityMenu}
            </Dropdown>
        </span>
    );
};
