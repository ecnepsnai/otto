import * as React from 'react';
import { Script, ScriptType } from '../../types/Script';
import { Link } from 'react-router-dom';
import { Menu } from '../../components/Menu';
import { Icon } from '../../components/Icon';
import { Table } from '../../components/Table';
import { GlobalModalFrame } from '../../components/Modal';
import { Rand } from '../../services/Rand';
import { RunModal } from '../run/RunModal';
import { Formatter } from '../../services/Formatter';

interface ScriptListItemProps {
    script: ScriptType;
    onReload: () => (void);
}
export const ScriptListItem: React.FC<ScriptListItemProps> = (props: ScriptListItemProps) => {
    const deleteMenuClick = () => {
        Script.DeleteModal(props.script).then(confirmed => {
            if (confirmed) {
                props.onReload();
            }
        });
    };

    const toggleMenuClick = () => {
        Script.Update(props.script, {
            Enabled: !props.script.Enabled
        }).then(() => {
            props.onReload();
        });
    };

    const enableDisableMenu = () => {
        if (props.script.Enabled) {
            return (<Menu.Item icon={<Icon.TimesCircle />} onClick={toggleMenuClick} label="Disable" />);
        }
        return (<Menu.Item icon={<Icon.CheckCircle />} onClick={toggleMenuClick} label="Enable" />);
    };

    const executeScriptMenuClick = () => {
        GlobalModalFrame.showModal(<RunModal scriptID={props.script.ID} key={Rand.ID()} />);
    };

    const link = <Link to={'/scripts/script/' + props.script.ID}>{props.script.Name}</Link>;

    return (
        <Table.Row disabled={!props.script.Enabled}>
            <td>{link}</td>
            <td>{props.script.Executable}</td>
            <td>{Formatter.ValueOrNothing((props.script.AttachmentIDs || []).length)}</td>
            <Table.Menu>
                <Menu.Item label="Run Script" icon={<Icon.PlayCircle />} onClick={executeScriptMenuClick} />
                <Menu.Link label="Edit" icon={<Icon.Edit />} to={'/scripts/script/' + props.script.ID + '/edit'} />
                {enableDisableMenu()}
                <Menu.Divider />
                <Menu.Item label="Delete" icon={<Icon.Delete />} onClick={deleteMenuClick} />
            </Table.Menu>
        </Table.Row>
    );
};
