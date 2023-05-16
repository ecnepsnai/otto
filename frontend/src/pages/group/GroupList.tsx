import * as React from 'react';
import { Group, GroupType } from '../../types/Group';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { CreateButton } from '../../components/Button';
import { Column, Table } from '../../components/Table';
import { Link } from 'react-router-dom';
import { ContextMenuItem } from '../../components/ContextMenu';
import { Icon } from '../../components/Icon';
import { Formatter } from '../../services/Formatter';
import { DefaultSort } from '../../services/Sort';
import { Permissions, UserAction } from '../../services/Permissions';

export const GroupList: React.FC = () => {
    const [loading, setLoading] = React.useState(true);
    const [groups, setGroups] = React.useState<GroupType[]>();
    const [membership, setMembership] = React.useState<{ [id: string]: string[] }>({});

    React.useEffect(() => {
        loadData();
    }, []);

    const loadGroups = () => {
        return Group.List().then(groups => {
            setGroups(groups);
        });
    };

    const loadMembership = () => {
        return Group.Membership().then(membership => {
            setMembership(membership);
        });
    };

    const loadData = () => {
        Promise.all([loadGroups(), loadMembership()]).then(() => {
            setLoading(false);
        });
    };

    if (loading) {
        return (<PageLoading />);
    }

    const toolbar = (
        <React.Fragment>
            <CreateButton to="/groups/group/" disabled={!Permissions.UserCan(UserAction.ModifyGroups)} />
        </React.Fragment>
    );

    const tableCols: Column[] = [
        {
            title: 'Name',
            value: (v: GroupType) => {
                return (<Link to={'/groups/group/' + v.ID}>{v.Name}</Link>);
            },
            sort: 'Name'
        },
        {
            title: 'Hosts',
            value: (v: GroupType) => {
                return (<span>{Formatter.ValueOrNothing(membership[v.ID].length)}</span>);
            },
            sort: (asc: boolean, left: GroupType, right: GroupType) => {
                return DefaultSort(asc, membership[left.ID].length, membership[right.ID].length);
            }
        },
        {
            title: 'Scripts',
            value: (v: GroupType) => {
                return (<span>{Formatter.ValueOrNothing((v.ScriptIDs || []).length)}</span>);
            },
            sort: (asc: boolean, left: GroupType, right: GroupType) => {
                return DefaultSort(asc, (left.ScriptIDs || []).length, (right.ScriptIDs || []).length);
            }
        },
    ];

    return (
        <Page title="Groups" toolbar={toolbar}>
            <Table columns={tableCols} data={groups} contextMenu={(a: GroupType) => GroupTableContextMenu(a, loadGroups)} defaultSort={{ ColumnIdx: 0, Ascending: true }} />
        </Page>
    );
};

const GroupTableContextMenu = (group: GroupType, onReload: () => (void)): (ContextMenuItem | 'separator')[] => {
    const deleteMenuClick = () => {
        Group.DeleteModal(group).then(confirmed => {
            if (confirmed) {
                onReload();
            }
        });
    };

    return [
        {
            title: 'Edit',
            icon: (<Icon.Edit />),
            href: '/groups/group/' + group.ID + '/edit',
            disabled: !Permissions.UserCan(UserAction.ModifyGroups),
        },
        'separator',
        {
            title: 'Delete',
            icon: (<Icon.Delete />),
            onClick: deleteMenuClick,
            disabled: !Permissions.UserCan(UserAction.ModifyGroups),
        },
    ];
};
