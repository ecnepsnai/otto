export class Clipboard {
    public static setText(text: string): Promise<void> {
        return new Promise(resolve => {
            navigator.clipboard.writeText(text).then(() => {
                resolve();
            }, () => {
                console.warn('Error setting clipboard contents - falling back to prompt');
                prompt('Copy the following text', text);
                resolve();
            });
        });
    }
}