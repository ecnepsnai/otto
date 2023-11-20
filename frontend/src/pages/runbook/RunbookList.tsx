import * as React from 'react';
import { Runbook, RunbookType } from '../../types/Runbook';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { CreateButton } from '../../components/Button';
import { Column, Table } from '../../components/Table';
import { Link } from 'react-router-dom';
import { ContextMenuItem } from '../../components/ContextMenu';
import { DateLabel } from '../../components/DateLabel';
import { Icon } from '../../components/Icon';
import { Permissions, UserAction } from '../../services/Permissions';

export const RunbookList: React.FC = () => {
    const [IsLoading, SetIsLoading] = React.useState(true);
    const [Runbooks, SetRunbooks] = React.useState<RunbookType[]>();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadRunbooks = () => {
        return Runbook.List().then(runbooks => {
            SetRunbooks(runbooks);
        });
    };

    const loadData = () => {
        Promise.all([loadRunbooks()]).then(() => {
            SetIsLoading(false);
        });
    };

    if (IsLoading) {
        return (<PageLoading />);
    }

    const toolbar = (
        <React.Fragment>
            <CreateButton to="/runbooks/runbook/" disabled={!Permissions.UserCan(UserAction.ModifyRunbooks)} />
        </React.Fragment>
    );

    const tableCols: Column[] = [
        {
            title: 'Name',
            value: (v: RunbookType) => {
                return (<Link to={'/runbooks/runbook/' + v.ID}>{v.Name}</Link>);
            },
            sort: 'Name'
        },
        {
            title: 'Last Run',
            value: (v: RunbookType) => {
                return (<DateLabel date={v.LastRun} />);
            }
        }
    ];

    return (
        <Page title="Runbooks" toolbar={toolbar}>
            <Table columns={tableCols} data={Runbooks} contextMenu={(a: RunbookType) => RunbookTableContextMenu(a, loadRunbooks)} defaultSort={{ ColumnIdx: 0, Ascending: true }} />
        </Page>
    );
};

const RunbookTableContextMenu = (runbook: RunbookType, onReload: () => void): (ContextMenuItem | 'separator')[] => {
    const deleteMenuClick = () => {
        Runbook.DeleteModal(runbook).then(confirmed => {
            if (confirmed) {
                onReload();
            }
        });
    };

    return [
        {
            title: 'Edit',
            icon: (<Icon.Edit />),
            href: '/runbooks/runbook/' + runbook.ID + '/edit',
            disabled: !Permissions.UserCan(UserAction.ModifyRunbooks),
        },
        'separator',
        {
            title: 'Delete',
            icon: (<Icon.Delete />),
            onClick: deleteMenuClick,
            disabled: !Permissions.UserCan(UserAction.ModifyRunbooks),
        },
    ];
};
