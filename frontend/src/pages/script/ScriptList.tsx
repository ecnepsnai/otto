import * as React from 'react';
import { Script, ScriptType } from '../../types/Script';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { CreateButton } from '../../components/Button';
import { Column, Table } from '../../components/Table';
import { Link } from 'react-router-dom';
import { ContextMenuItem } from '../../components/ContextMenu';
import { Icon } from '../../components/Icon';
import { GlobalModalFrame } from '../../components/Modal';
import { Formatter } from '../../services/Formatter';
import { Rand } from '../../services/Rand';
import { DefaultSort } from '../../services/Sort';
import { RunModal } from '../run/RunModal';

export const ScriptList: React.FC = () => {
    const [loading, setLoading] = React.useState<boolean>(true);
    const [scripts, setScripts] = React.useState<ScriptType[]>([]);

    React.useEffect(() => {
        loadScripts();
    }, []);

    const loadScripts = () => {
        Script.List().then(scripts => {
            setLoading(false);
            setScripts(scripts);
        });
    };

    if (loading) {
        return (<PageLoading />);
    }

    const toolbar = (
        <React.Fragment>
            <CreateButton to="/scripts/script/" />
        </React.Fragment>
    );

    const tableCols: Column[] = [
        {
            title: 'Name',
            value: (v: ScriptType) => {
                return (<Link to={'/scripts/script/' + v.ID}>{v.Name}</Link>);
            },
            sort: 'Name'
        },
        {
            title: 'Executable',
            value: 'Executable',
            sort: 'Executable'
        },
        {
            title: 'Attachments',
            value: (v: ScriptType) => {
                return (<span>{Formatter.ValueOrNothing((v.AttachmentIDs || []).length)}</span>);
            },
            sort: (asc: boolean, left: ScriptType, right: ScriptType) => {
                return DefaultSort(asc, (left.AttachmentIDs || []).length, (right.AttachmentIDs || []).length);
            }
        },
    ];

    return (
        <Page title="Scripts" toolbar={toolbar}>
            <Table columns={tableCols} data={scripts} contextMenu={(a: ScriptType) => ScriptTableContextMenu(a, loadScripts)} defaultSort={{ ColumnIdx: 0, Ascending: true }} />
        </Page>
    );
};

const ScriptTableContextMenu = (script: ScriptType, onReload: () => void): (ContextMenuItem | 'separator')[] => {
    const deleteMenuClick = () => {
        Script.DeleteModal(script).then(confirmed => {
            if (confirmed) {
                onReload();
            }
        });
    };

    const toggleMenuClick = () => {
        const s = script;
        s.Enabled = !s.Enabled;
        Script.Save(s).then(() => {
            onReload();
        });
    };

    const executeScriptMenuClick = () => {
        GlobalModalFrame.showModal(<RunModal scriptID={script.ID} key={Rand.ID()} />);
    };

    return [
        {
            title: 'Run Script',
            icon: (<Icon.PlayCircle />),
            onClick: executeScriptMenuClick,
        },
        {
            title: 'Edit',
            icon: (<Icon.Edit />),
            href: '/scripts/script/' + script.ID + '/edit',
        },
        {
            title: script.Enabled ? 'Disable' : 'Enable',
            icon: script.Enabled ? (<Icon.TimesCircle />) : (<Icon.CheckCircle />),
            onClick: toggleMenuClick,
        },
        'separator',
        {
            title: 'Delete',
            icon: (<Icon.Delete />),
            onClick: deleteMenuClick,
        },
    ];
};
