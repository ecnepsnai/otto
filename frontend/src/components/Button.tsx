import * as React from 'react';
import '../../css/button.scss';
import { Style } from './Style';
import { useNavigate } from 'react-router-dom';
import { Icon } from './Icon';

interface CommonButtonProps {
    onClick: () => (void);
    disabled?: boolean;
}

interface CommonButtonLinkProps {
    to: string;
    disabled?: boolean;
}

interface ButtonProps {
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
    onClick?: () => void;
    /**
     * Should the button be disabled
     */
    disabled?: boolean;
    /**
     * Optional classes to append to the button
     */
    className?: string;

    children?: React.ReactNode;
}

/**
 * A button that the user can click
 */
export const Button: React.FC<ButtonProps> = (props: ButtonProps) => {
    const onClick = () => {
        props.onClick();
    };

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
    return (
        <button type="button" className={className} onClick={onClick} disabled={props.disabled}>{props.children}</button>
    );
};

interface ButtonLinkProps {
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

    children?: React.ReactNode;
}

/**
 * A button that acts as a link. When clicked it redirects the user to the destination.
 */
export const ButtonLink: React.FC<ButtonLinkProps> = (props: ButtonLinkProps) => {
    const navigate = useNavigate();

    const onClick = () => {
        navigate(props.to);
    };
    return <Button
        color={props.color}
        outline={props.outline}
        size={props.size}
        disabled={props.disabled}
        onClick={onClick}>
        {props.children}
    </Button>;
};

interface ButtonsProps {
    children?: React.ReactNode
}
export const Buttons: React.FC<ButtonsProps> = (props: ButtonsProps) => {
    return (<div className="buttons">{props.children}</div>);
};

interface ButtonGroupProps {
    children?: React.ReactNode
}
export const ButtonGroup: React.FC<ButtonGroupProps> = (props: ButtonGroupProps) => {
    return (<div className="btn-group">{props.children}</div>);
};

export const AddButton: React.FC<CommonButtonProps> = (props: CommonButtonProps) => {
    return (<Button onClick={props.onClick} color={Style.Palette.Primary} outline size={Style.Size.S} disabled={props.disabled}>
        <Icon.Label label="Add" icon={<Icon.Plus />} />
    </Button>);
};

export const CreateButton: React.FC<CommonButtonLinkProps> = (props: CommonButtonLinkProps) => {
    return (<ButtonLink to={props.to} color={Style.Palette.Primary} outline size={Style.Size.S} disabled={props.disabled}>
        <Icon.Label label="Create New" icon={<Icon.Plus />} />
    </ButtonLink>);
};

export const EditButton: React.FC<CommonButtonLinkProps> = (props: CommonButtonLinkProps) => {
    return (<ButtonLink to={props.to} color={Style.Palette.Primary} outline size={Style.Size.S} disabled={props.disabled}>
        <Icon.Label label="Edit" icon={<Icon.Edit />} />
    </ButtonLink>);
};

export const DeleteButton: React.FC<CommonButtonProps> = (props: CommonButtonProps) => {
    return (<Button onClick={props.onClick} color={Style.Palette.Danger} outline size={Style.Size.S} disabled={props.disabled}>
        <Icon.Label label="Delete" icon={<Icon.Delete />} />
    </Button>);
};

export const SmallPlayButton: React.FC<CommonButtonProps> = (props: CommonButtonProps) => {
    return (
        <button className="play-button" onClick={props.onClick}>
            <Icon.PlayCircle />
        </button>
    );
};

interface ButtonAnchorProps {
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
    href: string;
    /**
     * Target value for the anchor
     */
    target?: string;
    /**
     * Should the destination be downloaded
     */
    download?: boolean;
    children?: React.ReactNode;
}

/**
 * A true anchor button. Should only be used when actual navigation is required, otherwise use ButtonLink.
 */
export const ButtonAnchor: React.FC<ButtonAnchorProps> = (props: ButtonAnchorProps) => {
    let className = 'btn ';
    if (props.outline) {
        className += 'btn-outline-';
    } else {
        className += 'btn-';
    }
    className += props.color.toString();
    const size = props.size || Style.Size.S;
    className += ' btn-' + size.toString();
    return (<a href={props.href} download={props.download} className={className} target={props.target}>{props.children}</a>);
};

/**
 * A button where the user must click twice within 5 seconds before the action is performed.
 */
export const ConfirmButton: React.FC<ButtonProps> = (props: ButtonProps) => {
    const [didClick, setDidClick] = React.useState(false);

    React.useEffect(() => {
        if (!didClick) {
            return;
        }

        setTimeout(() => {
            if (didClick) {
                setDidClick(false);
            }
        }, 5000);
    }, [didClick]);

    const onClick = () => {
        if (didClick) {
            setDidClick(false);
            props.onClick();
        } else {
            setDidClick(true);
        }
    };

    if (didClick) {
        return (<Button
            color={props.color}
            outline={props.outline}
            size={props.size}
            onClick={onClick}
            disabled={props.disabled}
            className={props.className}>Confirm?</Button>);
    }

    return (<Button
        color={props.color}
        outline={props.outline}
        size={props.size}
        onClick={onClick}
        disabled={props.disabled}
        className={props.className}>{props.children}</Button>);
};
