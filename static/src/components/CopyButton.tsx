import * as React from 'react';
import { Clipboard } from '../services/Clipboard';
import { Button } from './Button';
import { Icon } from './Icon';
import { Style } from './Style';

export interface CopyButtonProps {
    text: string;
}
interface CopyButtonState {
    didCopy: boolean;
}
export class CopyButton extends React.Component<CopyButtonProps, CopyButtonState> {
    constructor(props: CopyButtonProps) {
        super(props);
        this.state = {
            didCopy: false,
        };
    }

    private onClick = () => {
        Clipboard.setText(this.props.text).then(() => {
            this.setState({ didCopy: true });
        });
    }

    private content = () => {
        if (this.state.didCopy) {
            return (<Icon.CheckCircle />);
        }

        return (<Icon.Clipboard />);
    }

    render(): JSX.Element {
        const color = this.state.didCopy ? Style.Palette.Success : Style.Palette.Primary;
        return (
            <Button color={color} outline size={Style.Size.XS} onClick={this.onClick} disabled={this.state.didCopy}>
                { this.content() }
            </Button>
        );
    }
}
