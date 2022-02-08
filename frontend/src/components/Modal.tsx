import * as React from 'react';
import { Style } from './Style';
import { Button } from './Button';
import { Rand } from '../services/Rand';
import { Form } from './Form';
import { Modal as BSModal } from 'bootstrap';

export interface ModalButton {
    /**
     * The label for the button
     */
    label: JSX.Element | string;
    /**
     * The color of the button
     */
    color?: Style.Palette;
    /**
     * Event fired when the button is clicked
     */
    onClick?: () => void;
    /**
     * If true, clicking this button will not dismiss the modal
     */
    dontDismiss?: boolean;
    /**
     * If true the button is disabled
     */
    disabled?: boolean;
}

interface ModalProps {
    /**
     * The title of the modal
     */
    title?: string;
    /**
     * The header of the modal
     */
    header?: JSX.Element;
    /**
     * Array of buttons for the modal
     */
    buttons?: ModalButton[];
    /**
     * Event fired after the modal was dismissed and is no longer visible
     */
    dismissed?: () => void;
    /**
     * Optional size for the modal
     */
    size?: Style.Size;
    /**
     * Optional if the modal can not be dismissed by clicking the background or pressing esc
     */
    static?: boolean;
}

interface ModalState {
    id: string;
    bsModal?: BSModal;
}

export class Modal extends React.Component<ModalProps, ModalState> {
    private static currentModal: BSModal = undefined;
    constructor(props: ModalProps) {
        super(props);
        this.state = { id: Rand.ID() };
    }
    componentDidMount(): void {
        const element = document.getElementById(this.state.id);
        element.addEventListener('hidden.bs.modal', () => {
            if (this.props.dismissed) {
                this.props.dismissed();
            }
            GlobalModalFrame.removeModal();
            Modal.currentModal = undefined;
        });
        let backdrop: 'static' | boolean = true;
        if (this.props.static) {
            backdrop = 'static';
        }
        const bsm = new BSModal(element, { backdrop: backdrop });
        bsm.show();
        this.setState({ bsModal: bsm });
        Modal.currentModal = bsm;
    }
    private buttonClick = (button: ModalButton) => {
        return () => {
            if (button.onClick) {
                button.onClick();
            }
            if (!button.dontDismiss) {
                this.state.bsModal.hide();
            }
        };
    };
    private closeButtonClick = () => {
        this.state.bsModal.hide();
    };
    private closeButton = () => {
        if (this.props.static) {
            return null;
        }
        return (
            <button type="button" onClick={this.closeButtonClick} className="btn-close" data-dismiss="modal" aria-label="Close"></button>
        );
    };
    private header = () => {
        if (this.props.title) {
            return (
                <div className="modal-header">
                    <h5 className="modal-title">{this.props.title}</h5>
                    { this.closeButton()}
                </div>
            );
        } else if (this.props.header) {
            return this.props.header;
        }
        return null;
    };
    private footer = () => {
        if (!this.props.buttons || this.props.buttons.length == 0) {
            return null;
        }

        return (
            <div className="modal-footer">
                {
                    this.props.buttons.map(button => {
                        button.color = button.color ?? Style.Palette.Primary;
                        return <Button color={button.color} onClick={this.buttonClick(button)} key={Rand.ID()} disabled={button.disabled}>{button.label}</Button>;
                    })
                }
            </div>
        );
    };
    render(): JSX.Element {
        let className = 'modal-dialog';
        if (this.props.size) {
            className += ' modal-' + this.props.size.toString();
        }

        return (
            <div className="modal fade" id={this.state.id}>
                <div className={className}>
                    <div className="modal-content">
                        {this.header()}
                        <div className="modal-body">
                            {this.props.children}
                        </div>
                        {this.footer()}
                    </div>
                </div>
            </div>
        );
    }

    public static dismiss(): void {
        if (!Modal.currentModal) {
            console.warn('Cannot dismiss modal when none present');
            return;
        }
        Modal.currentModal.hide();
    }

    /**
     * A confirm dialog where the two buttons are 'Cancel' and 'Delete'
     * @param title The title of the dialog
     * @param body The body of the dialog
     * @returns A promise that is resolved with wether or not the user clicked the 'Delete' button
     */
    public static delete(title: string, body: string): Promise<boolean> {
        return new Promise(resolve => {
            const buttonClick = (confirm: boolean): () => (void) => {
                return () => {
                    resolve(confirm);
                };
            };
            const dismissed = () => {
                buttonClick(false)();
            };
            const buttons: ModalButton[] = [
                {
                    label: 'Cancel',
                    color: Style.Palette.Secondary,
                    onClick: buttonClick(false),
                },
                {
                    label: 'Delete',
                    color: Style.Palette.Danger,
                    onClick: buttonClick(true),
                }
            ];
            GlobalModalFrame.showModal(
                <Modal title={title} buttons={buttons} dismissed={dismissed}>
                    <p>{body}</p>
                </Modal>
            );
        });
    }

    /**
     * A confirm dialog where the two buttons are 'Cancel' and 'Confirm'. Not sutible for dangerous actions.
     * @param title The title of the dialog
     * @param body The body of the dialog
     * @returns A promise that is resolved with wether or not the user clicked the 'Confirm' button
     */
    public static confirm(title: string, body: string|JSX.Element): Promise<boolean> {
        return new Promise(resolve => {
            const buttonClick = (confirm: boolean): () => (void) => {
                return () => {
                    resolve(confirm);
                };
            };
            const dismissed = () => {
                buttonClick(false)();
            };
            const buttons: ModalButton[] = [
                {
                    label: 'Cancel',
                    color: Style.Palette.Secondary,
                    onClick: buttonClick(false),
                },
                {
                    label: 'Confirm',
                    color: Style.Palette.Primary,
                    onClick: buttonClick(true),
                }
            ];
            GlobalModalFrame.showModal(
                <Modal title={title} buttons={buttons} dismissed={dismissed}>
                    {typeof body === 'string' ? (<p>{body}</p>) : body}
                </Modal>
            );
        });
    }
}

interface ModalHeaderProps {
    children?: React.ReactNode
}
export const ModalHeader: React.FC<ModalHeaderProps> = (props: ModalHeaderProps) => {
    const closeButtonClicked = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
        Modal.dismiss();
    };

    return (
        <div className="modal-header">
            { props.children}
            <button type="button" className="btn-close" data-dismiss="modal" aria-label="Close" onClick={closeButtonClicked}></button>
        </div>
    );
};

interface GlobalModalFrameState {
    modal?: JSX.Element;
}

export class GlobalModalFrame extends React.Component<unknown, GlobalModalFrameState> {
    constructor(props: unknown) {
        super(props);
        this.state = {};
        GlobalModalFrame.instance = this;
    }

    private static instance: GlobalModalFrame;

    public static showModal(modal: JSX.Element): void {
        this.instance.setState(state => {
            if (state.modal != undefined) {
                throw new Error('Refusing to stack modals');
            }
            return { modal: modal };
        });
    }

    public static removeModal(): void {
        try {
            document.body.classList.remove('modal-open');
            document.body.removeAttribute('style');
            document.querySelector('.modal-backdrop').remove();
        } catch (e) {
            //
        }
        this.instance.setState({ modal: undefined });
    }

    render(): JSX.Element {
        return (
            <div id="global-modal-frame">
                {
                    this.state.modal
                }
            </div>
        );
    }
}

interface ModalFormProps {
    /**
     * The title of the modal
     */
    title: string;
    /**
     * Submit method - called when the submit button is clicked if the form is valid. Return a promise that when resolved
     * will dismiss the modal.
     */
    onSubmit: () => (Promise<unknown>);
    /**
     * Called when the modal is dismissed without saving
     */
    onDismissed?: () => void;

    children?: React.ReactNode;
}
/**
 * A modal with a form element. Buttons in the footer of the modal are used for the form element.
 */
export const ModalForm: React.FC<ModalFormProps> = (props: ModalFormProps) => {
    const [loading, setLoading] = React.useState(false);
    const formRef: React.RefObject<Form> = React.createRef();

    const saveClick = () => {
        if (!formRef.current.validateForm()) {
            return;
        }

        setLoading(true);
        props.onSubmit().then(() => {
            GlobalModalFrame.removeModal();
        }, err => {
            console.error('Error while executing modal form save promise', err);
            setLoading(false);
        });
    };

    const buttons: ModalButton[] = [
        {
            label: 'Discard',
            color: Style.Palette.Secondary,
            disabled: loading,
            onClick: props.onDismissed,
        },
        {
            label: 'Save',
            color: Style.Palette.Primary,
            onClick: saveClick,
            dontDismiss: true,
            disabled: loading,
        }
    ];

    return (
        <Modal title={props.title} buttons={buttons} static>
            <Form ref={formRef}>
                {props.children}
            </Form>
        </Modal>
    );
};
