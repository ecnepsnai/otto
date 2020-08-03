/* eslint-disable @typescript-eslint/no-explicit-any */
export interface ToastOptions { animation?: boolean; autohide?: boolean; delay?: number; }
export interface ModalOptions { backdrop?: boolean | string; keyboard?: boolean; focus?: boolean; show?: boolean; hide?: boolean; }

export type BSModule = any;

export class Bootstrap {
    private static bs: any = ((window as any).bootstrap as any);

    public static Toast(element: HTMLElement, options: ToastOptions): BSModule {
        return new this.bs.Toast(element, options);
    }

    public static Alert(element: HTMLElement): BSModule {
        return new this.bs.Alert(element);
    }

    public static Modal(element: HTMLElement, options: ModalOptions): BSModule {
        return new this.bs.Modal(element, options);
    }

    public static Popover(element: HTMLElement): BSModule {
        return new this.bs.Popover(element);
    }
}