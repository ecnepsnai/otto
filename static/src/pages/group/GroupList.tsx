import * as React from 'react';
import { Group, GroupType } from '../../types/Group';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Buttons, CreateButton } from '../../components/Button';
import { Table } from '../../components/Table';
import { GroupListItem } from './GroupListItem';

interface GroupListState {
    loading: boolean;
    groups: GroupType[];
    membership: {[id: string]: string[]};
}
export class GroupList extends React.Component<unknown, GroupListState> {
    constructor(props: unknown) {
        super(props);
        this.state = {
            loading: true,
            groups: [],
            membership: {},
        };
    }

    componentDidMount(): void {
        this.loadData();
    }

    private loadGroups = () => {
        return Group.List().then(groups => {
            this.setState({
                groups: groups,
            });
        });
    }

    private loadMembership = () => {
        return Group.Membership().then(membership => {
            this.setState({
                membership: membership,
            });
        });
    }

    private loadData = () => {
        Promise.all([this.loadGroups(), this.loadMembership()]).then(() => {
            this.setState({ loading: false });
        });
    }

    render(): JSX.Element {
        if (this.state.loading) {
            return ( <PageLoading /> );
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
                            this.state.groups.map(group => {
                                return <GroupListItem group={group} hosts={this.state.membership[group.ID]} key={group.ID} onReload={this.loadData} numGroups={this.state.groups.length}></GroupListItem>;
                            })
                        }
                    </Table.Body>
                </Table.Table>
            </Page>
        );
    }
}
