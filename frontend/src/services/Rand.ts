/**
 * Class for generating random things
 */
export class Rand {
    private static alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';

    /**
     * Generate a mathematically random ID string
     */
    public static ID(): string {
        let text = '';
        for (let i = 0; i < 10; i++) {
            text += Rand.alphabet.charAt(Math.floor(Math.random() * Rand.alphabet.length));
        }
        return text;
    }

    /**
     * Generate a cryptographically suitable random string
     */
    public static PSK(): string {
        const array = new Uint32Array(6);
        window.crypto.getRandomValues(array);
        return Array.from(array, byte => byte.toString(16).padStart(2, '0')).join('');
    }
}
