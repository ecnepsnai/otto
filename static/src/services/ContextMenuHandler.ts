import * as React from 'react';

/**
 * Class for handling context menu events
 */
export class ContextMenuHandler {
    private callback: (x: number, y: number) => void;
    private longPressTimeout: NodeJS.Timeout;
    private readonly longPressDurationMS = 610;

    /**
     * Create a new instance of the menu handler. Each instance should only be associated with one single react element.
     * @param callback Called when the context menu should be shown
     */
    constructor(callback: (x: number, y: number) => void) {
        this.callback = callback;
    }

    /**
     * Event handler for onContextMenu
     */
    public readonly onContextMenu = (event: React.MouseEvent): void => {
        event.preventDefault();
        this.callback(event.pageX, event.pageY);
    };

    /**
     * Event handler for onTouchStart
     */
    public readonly onTouchStart = (event: React.TouchEvent): void => {
        // Ignore touches on anchors
        if ((event.touches[0].target as HTMLElement).nodeName === 'A') {
            return;
        }

        this.longPressTimeout = setTimeout(() => {
            this.callback(event.touches[0].pageX, event.touches[0].pageY);
        }, this.longPressDurationMS);
    };

    /**
     * Event handler for onTouchEnd
     */
    public readonly onTouchEnd = (): void => {
        clearTimeout(this.longPressTimeout);
    };

    /**
     * Event handler for onTouchMove
     */
    public readonly onTouchMove = (): void => {
        clearTimeout(this.longPressTimeout);
    };
}