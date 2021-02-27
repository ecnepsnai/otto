import * as React from 'react';
import { Clipboard } from '../services/Clipboard';
import { Button } from './Button';
import { Icon } from './Icon';
import { Style } from './Style';

export interface CopyButtonProps {
    text: string;
}
export const CopyButton: React.FC<CopyButtonProps> = (props: CopyButtonProps) => {
    const [didCopy, setDidCopy] = React.useState(false);

    const onClick = () => {
        Clipboard.setText(props.text).then(() => {
            setDidCopy(true);
        });
    };

    const content = () => {
        if (didCopy) {
            return (<Icon.CheckCircle />);
        }

        return (<Icon.Clipboard />);
    };

    const color = didCopy ? Style.Palette.Success : Style.Palette.Primary;
    return (
        <Button color={color} outline size={Style.Size.XS} onClick={onClick} disabled={didCopy}>
            { content() }
        </Button>
    );
};
