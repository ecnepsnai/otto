import * as React from 'react';
import { Group, GroupType } from '../../types/Group';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Buttons, CreateButton } from '../../components/Button';
import { Table } from '../../components/Table';
import { GroupListItem } from './GroupListItem';

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

    return (
        <Page title="Groups">
            <Buttons>
                <CreateButton to="/groups/group/" />
            </Buttons>
            <Table.Table>
                <Table.Head>
                    <Table.Column>Name</Table.Column>
                    <Table.Column>Hosts</Table.Column>
                    <Table.Column>Scripts</Table.Column>
                    <Table.MenuColumn />
                </Table.Head>
                <Table.Body>
                    {
                        groups.map(group => {
                            return <GroupListItem group={group} hosts={membership[group.ID]} key={group.ID} onReload={loadData} numGroups={groups.length}></GroupListItem>;
                        })
                    }
                </Table.Body>
            </Table.Table>
        </Page>
    );
};
