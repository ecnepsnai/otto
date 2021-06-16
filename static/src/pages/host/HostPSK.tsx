import * as React from 'react';
import { Host, HostType } from '../../types/Host';
import { ListGroup } from '../../components/ListGroup';
import { Dropdown, Menu } from '../../components/Menu';
import { Icon } from '../../components/Icon';
import { Clipboard } from '../../services/Clipboard';
import { Notification } from '../../components/Notification';

interface HostPSKProps {
    host: HostType;
    didRotate: (newPSK: string) => void;
}
export const HostPSK: React.FC<HostPSKProps> = (props: HostPSKProps) => {
    const [IsLoading, setIsLoading] = React.useState(false);

    const copyClick = () => {
        Clipboard.setText(props.host.PSK);
    };

    const rotateClick = () => {
        setIsLoading(true);
        Host.RotatePSK(props.host.ID).then(newPSK => {
            setIsLoading(false);
            Notification.success('Host PSK rotated');
            props.didRotate(newPSK);
        });
    };

    const content = () => {
        if (IsLoading) {
            return (<Icon.Spinner pulse />);
        }

        return (<Dropdown label={<Icon.Bars />}>
            <Menu.Item icon={<Icon.Clipboard />} label="Copy" onClick={copyClick} />
            <Menu.Item icon={<Icon.Random />} label="Rotate" onClick={rotateClick} />
        </Dropdown>);
    };

    return (
        <ListGroup.TextItem title="PSK">
            <code>*****</code>
            <span className="ms-2">{content()}</span>
        </ListGroup.TextItem>
    );
};
