import * as React from 'react';
import { Script, ScriptType } from '../../types/Script';
import { Link } from 'react-router-dom';
import { Icon } from '../../components/Icon';
import { Table } from '../../components/Table';
import { GlobalModalFrame } from '../../components/Modal';
import { Rand } from '../../services/Rand';
import { RunModal } from '../run/RunModal';
import { Formatter } from '../../services/Formatter';
import { ContextMenuItem } from '../../components/ContextMenu';

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
        const s = props.script;
        s.Enabled = !s.Enabled;
        Script.Save(s).then(() => {
            props.onReload();
        });
    };

    const executeScriptMenuClick = () => {
        GlobalModalFrame.showModal(<RunModal scriptID={props.script.ID} key={Rand.ID()} />);
    };

    const link = <Link to={'/scripts/script/' + props.script.ID}>{props.script.Name}</Link>;

    const contextMenu: (ContextMenuItem | 'separator')[] = [
        {
            title: 'Run Script',
            icon: (<Icon.PlayCircle />),
            onClick: executeScriptMenuClick,
        },
        {
            title: 'Edit',
            icon: (<Icon.Edit />),
            href: '/scripts/script/' + props.script.ID + '/edit',
        },
        {
            title: props.script.Enabled ? 'Disable' : 'Enable',
            icon: props.script.Enabled ? (<Icon.TimesCircle />) : (<Icon.CheckCircle />),
            onClick: toggleMenuClick,
        },
        'separator',
        {
            title: 'Delete',
            icon: (<Icon.Delete />),
            onClick: deleteMenuClick,
        },
    ];

    return (
        <Table.Row disabled={!props.script.Enabled} menu={contextMenu}>
            <td>{link}</td>
            <td>{props.script.Executable}</td>
            <td>{Formatter.ValueOrNothing((props.script.AttachmentIDs || []).length)}</td>
        </Table.Row>
    );
};
