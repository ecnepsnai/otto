import * as React from 'react';
import { Schedule } from '../../types/Schedule';
import { Link } from 'react-router-dom';
import { Menu } from '../../components/Menu';
import { Icon } from '../../components/Icon';
import { Table } from '../../components/Table';
import { Script } from '../../types/Script';
import { EnabledBadge } from '../../components/Badge';
import { SchedulePattern } from './SchedulePattern';
import { DateLabel } from '../../components/DateLabel';

export interface ScheduleListItemProps { schedule: Schedule, script: Script, onReload: () => (void); }
export class ScheduleListItem extends React.Component<ScheduleListItemProps, {}> {
    private deleteMenuClick = () => {
        this.props.schedule.DeleteModal().then(confirmed => {
            if (confirmed) {
                this.props.onReload();
            }
        });
    }

    private enabledOnColumn = () => {
        if (this.props.schedule.Scope.GroupIDs.length > 0) {
            let unit = 'groups';
            if (this.props.schedule.Scope.GroupIDs.length  == 1) {
                unit = 'group';
            }
            return (<td>{this.props.schedule.Scope.GroupIDs.length} {unit}</td>);
        } else if (this.props.schedule.Scope.HostIDs.length > 0) {
            let unit = 'hosts';
            if (this.props.schedule.Scope.HostIDs.length == 1) {
                unit = 'host';
            }
            return (<td>{this.props.schedule.Scope.HostIDs.length} {unit}</td>);
        }

        return (<td></td>);
    }

    render(): JSX.Element {
        const link = <Link to={'/schedules/schedule/' + this.props.schedule.ID}>{ this.props.schedule.Name }</Link>;

        return (
            <Table.Row>
                <td>{ link }</td>
                <td>{ this.props.script.Name }</td>
                <td><SchedulePattern pattern={this.props.schedule.Pattern} /></td>
                {this.enabledOnColumn()}
                <td><DateLabel date={this.props.schedule.LastRunTime} /></td>
                <td><EnabledBadge value={this.props.schedule.Enabled} /></td>
                <Table.Menu>
                    <Menu.Link label="Edit" icon={<Icon.Edit />} to={'/schedules/schedule/' + this.props.schedule.ID + '/edit'}/>
                    <Menu.Divider />
                    <Menu.Item label="Delete" icon={<Icon.Delete />} onClick={this.deleteMenuClick}/>
                </Table.Menu>
            </Table.Row>
        );
    }
}
