export class Clipboard {
    private static promptCopy(text: string): Promise<void> {
        return new Promise(resolve => {
            prompt('Copy the following text', text);
            resolve();
        });
    }
    public static setText(text: string): Promise<void> {
        if (!navigator.clipboard) {
            return Clipboard.promptCopy(text);
        }

        return new Promise(resolve => {
            navigator.clipboard.writeText(text).then(() => {
                resolve();
            }, () => {
                return Clipboard.promptCopy(text);
            }).catch(() => {
                return Clipboard.promptCopy(text);
            });
        });
    }
}