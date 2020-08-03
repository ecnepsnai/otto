import * as React from 'react';

import { Host } from '../../types/Host';
import { Link } from 'react-router-dom';
import { Dropdown, MenuItem } from '../../components/Menu';
import { Style } from '../../components/Style';
import { Icon } from '../../components/Icon';
import { Table } from '../../components/Table';
import { HeartbeatBadge } from '../../components/Badge';
import { Heartbeat } from '../../types/Heartbeat';

export interface HostListItemProps { host: Host, heartbeat: Heartbeat, onReload: () => (void); }

interface HostListItemState { navigateToEdit?: boolean }

export class HostListItem extends React.Component<HostListItemProps, HostListItemState> {
    constructor(props: HostListItemProps) {
        super(props);
        this.state = { };
    }

    private editMenuClick = () => {
        this.setState({ navigateToEdit: true });
    }

    private deleteMenuClick = () => {
        this.props.host.DeleteModal().then(confirmed => {
            if (confirmed) {
                this.props.onReload();
            }
        });
    }

    render(): JSX.Element {
        const link = <Link to={'/hosts/host/' + this.props.host.ID}>{ this.props.host.Name }</Link>;
        const dropdownLabel = <Icon.Bars />;
        const buttonProps = {
            color: Style.Palette.Secondary,
            outline: true,
            size: Style.Size.XS,
        };
        return (
            <Table.Row disabled={!this.props.host.Enabled}>
                <td>{ link }</td>
                <td>{ this.props.host.Address }</td>
                <td><HeartbeatBadge heartbeat={this.props.heartbeat} /></td>
                <td>
                    <Dropdown label={dropdownLabel} button={buttonProps}>
                        <MenuItem label="Edit" icon={<Icon.Edit />} onClick={this.editMenuClick}/>
                        <MenuItem label="Delete" icon={<Icon.Delete />} onClick={this.deleteMenuClick}/>
                    </Dropdown>
                </td>
            </Table.Row>
        );
    }
}
