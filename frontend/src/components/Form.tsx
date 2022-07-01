import * as React from 'react';
import { Style } from './Style';
import { Button } from './Button';
import { Icon } from './Icon';
import '../../css/form.scss';

export interface ValidationResult {
    valid: boolean;
    invalidMessage?: string;
}

interface FormProps {
    /**
     * Should a save button be appended to the bottom of the form
     */
    showSaveButton?: boolean;
    /**
     * Event fired when the form is submitted
     */
    onSubmit?: () => void;
    /**
     * If the form should span the full width of the page
     */
    fullWidth?: boolean;
    /**
     * If true the submit button is disabled
     */
    loading?: boolean;

    children?: React.ReactNode;
}

interface FormState {
    invalid?: boolean;
}

/**
 * A form is a collection of inputs and a submit button
 */
export class Form extends React.Component<FormProps, FormState> {
    private domRef: React.RefObject<HTMLFormElement>;

    constructor(props: FormProps) {
        super(props);
        this.domRef = React.createRef();
        this.state = {};
    }

    /**
     * Performs validation on the form
     * @returns true if the form is valid, false if invalid
     */
    public validateForm = (): boolean => {
        const elemn = this.domRef.current;
        const invalidNodes = elemn.querySelectorAll('[data-valid="invalid"]');
        if (invalidNodes.length > 0) {
            this.setState({ invalid: true });
            return false;
        }
        return true;
    };

    private onClick = () => {
        this.submitForm();
    };

    private onSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        this.submitForm();
    };

    private submitForm = () => {
        this.validateForm();

        if (this.props.onSubmit) {
            this.props.onSubmit();
        }
    };

    private saveButton = () => {
        if (!this.props.showSaveButton) {
            return null;
        }

        let content = (<Icon.Label icon={<Icon.CheckCircle />} label="Apply" />);
        if (this.props.loading) {
            content = (<Icon.Label icon={<Icon.Spinner pulse />} label="Loading..." />);
        }

        return (
            <div className="mt-3">
                <Button color={Style.Palette.Primary} size={Style.Size.M} onClick={this.onClick} disabled={this.props.loading}>
                    {content}
                </Button>
            </div>
        );
    };

    private error = () => {
        if (!this.state.invalid) {
            return null;
        }

        return (
            <div className="mt-2">
                <Icon.Label icon={<Icon.TimesCircle color={Style.Palette.Danger} />} label="Correct Errors Before Continuing" />
            </div>
        );
    };

    render(): JSX.Element {
        const className = this.props.fullWidth ? '' : 'container';
        return (
            <form onSubmit={this.onSubmit} ref={this.domRef} className={className}>
                <fieldset>{this.props.children}</fieldset>
                {this.saveButton()}
                {this.error()}
            </form>
        );
    }
}

interface FormGroupProps {
    className?: string;
    children?: React.ReactNode;
    thin?: boolean;
}
export const FormGroup: React.FC<FormGroupProps> = (props: FormGroupProps) => {
    const margin = props.thin ? 'mb-1' : 'mb-3';
    return (<div className={(props.className ?? '') + ' ' + margin}>{props.children}</div>);
};
