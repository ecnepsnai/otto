import * as React from 'react';
import '../../css/button.scss';
import { Style } from './Style';
import { Redirect } from 'react-router-dom';
import { Icon } from './Icon';

export interface ButtonProps {
    /**
     * The color of this button
     */
    color: Style.Palette;
    /**
     * If an outline style should be used for this button
     */
    outline?: boolean;
    /**
     * Optional size for the button, defaults to regular
     */
    size?: Style.Size;
    /**
     * Event called when the button is clicked
     */
    onClick?: Function;
    /**
     * Should the button be disabled
     */
    disabled?: boolean;
    /**
     * Optional classes to append to the button
     */
    className?: string;
}

/**
 * A button that the user can click
 */
export class Button extends React.Component<ButtonProps, {}> {
    constructor(props: ButtonProps) {
        super(props);
    }
    private onClick = () => {
        this.props.onClick();
    }
    public static className(props: ButtonProps): string {
        let className = 'btn ';
        if (props.outline) {
            className += 'btn-outline-';
        } else {
            className += 'btn-';
        }
        className += props.color.toString();
        const size = props.size || Style.Size.S;
        className += ' btn-' + size.toString();
        className += ' ' + (props.className || '');
        return className;
    }
    render(): JSX.Element {
        const className = Button.className(this.props);
        return (
            <button type="button" className={className} onClick={this.onClick} disabled={this.props.disabled}>{this.props.children}</button>
        );
    }
}

export interface ButtonLinkProps {
    /**
     * The color of this button
     */
    color: Style.Palette;
    /**
     * If an outline style should be used for this button
     */
    outline?: boolean;
    /**
     * Optional size for the button, defaults to regular
     */
    size?: Style.Size;
    /**
     * The destination for the link
     */
    to: string;
    /**
     * Should the button be disabled
     */
    disabled?: boolean;
}
interface ButtonLinkState { doNavigate: boolean }

/**
 * A button that acts as a link. When clicked it redirects the user to the destination.
 */
export class ButtonLink extends React.Component<ButtonLinkProps, ButtonLinkState> {
    constructor(props: ButtonLinkProps) {
        super(props);
        this.state = { doNavigate: false };
    }
    private onClick = () => {
        this.setState({ doNavigate: true });
    }
    render(): JSX.Element {
        if (this.state.doNavigate) {
            return <Redirect push to={this.props.to} />;
        }
        return <Button
            color={this.props.color}
            outline={this.props.outline}
            size={this.props.size}
            disabled={this.props.disabled}
            onClick={this.onClick}>
                { this.props.children }
            </Button>;
    }
}

export class Buttons extends React.Component<{}, {}> {
    render(): JSX.Element {
        return <div className="buttons">{ this.props.children }</div>;
    }
}

export class ButtonGroup extends React.Component<{}, {}> {
    render(): JSX.Element {
        return <div className="btn-group">{ this.props.children }</div>;
    }
}

export interface CreateButtonProps {
    onClick?: () => (void);
    to?: string;
}

export class CreateButton extends React.Component<CreateButtonProps, {}> {
    render(): JSX.Element {
        if (this.props.onClick != undefined) {
            return (
                <Button onClick={this.props.onClick} color={Style.Palette.Primary} outline size={Style.Size.S}>
                    <Icon.Label label="Create New" icon={<Icon.Plus />} />
                </Button>
            );
        }

        if (this.props.to != undefined) {
            return (
                <ButtonLink to={this.props.to} color={Style.Palette.Primary} outline size={Style.Size.S}>
                    <Icon.Label label="Create New" icon={<Icon.Plus />} />
                </ButtonLink>
            );
        }

        return null;
    }
}

export interface EditButtonProps {
    to: string;
}

export class EditButton extends React.Component<EditButtonProps, {}> {
    render(): JSX.Element {
        return (
            <ButtonLink to={this.props.to} color={Style.Palette.Primary} outline size={Style.Size.S}>
                <Icon.Label label="Edit" icon={<Icon.Edit />} />
            </ButtonLink>
        );
    }
}

export interface DeleteButtonProps {
    onClick: () => (void);
    disabled?: boolean;
}

export class DeleteButton extends React.Component<DeleteButtonProps, {}> {
    render(): JSX.Element {
        return (
            <Button onClick={this.props.onClick} color={Style.Palette.Danger} outline size={Style.Size.S} disabled={this.props.disabled}>
                <Icon.Label label="Delete" icon={<Icon.Delete />} />
            </Button>
        );
    }
}

export interface SmallPlayButtonProps {
    onClick: () => (void);
}
export class SmallPlayButton extends React.Component<SmallPlayButtonProps, {}> {
    render(): JSX.Element {
        return (
            <button className="play-button" onClick={this.props.onClick}>
                <Icon.PlayCircle />
            </button>
        );
    }
}
