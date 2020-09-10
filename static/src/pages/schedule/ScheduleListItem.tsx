import * as React from 'react';
import { Schedule } from '../../types/Schedule';
import { Link } from 'react-router-dom';
import { Dropdown, MenuItem } from '../../components/Menu';
import { Style } from '../../components/Style';
import { Icon } from '../../components/Icon';
import { Table } from '../../components/Table';
import { Script } from '../../types/Script';
import { EnabledBadge } from '../../components/Badge';
import { SchedulePattern } from './SchedulePattern';
import { DateLabel } from '../../components/DateLabel';

export interface ScheduleListItemProps { schedule: Schedule, script: Script, onReload: () => (void); }

interface ScheduleListItemState { navigateToEdit?: boolean }

export class ScheduleListItem extends React.Component<ScheduleListItemProps, ScheduleListItemState> {
    constructor(props: ScheduleListItemProps) {
        super(props);
        this.state = { };
    }

    private editMenuClick = () => {
        this.setState({ navigateToEdit: true });
    }

    private deleteMenuClick = () => {
        this.props.schedule.DeleteModal().then(confirmed => {
            if (confirmed) {
                this.props.onReload();
            }
        });
    }

    private enabledOnColumn = () => {
        if (this.props.schedule.Scope.GroupIDs.length > 0) {
            return (<td>{this.props.schedule.Scope.GroupIDs.length} groups</td>);
        } else if (this.props.schedule.Scope.HostIDs.length > 0) {
            return (<td>{this.props.schedule.Scope.HostIDs.length} hosts</td>);
        }

        return (<td></td>);
    }

    render(): JSX.Element {
        const link = <Link to={'/schedules/schedule/' + this.props.schedule.ID}>{ this.props.script.Name }</Link>;
        const dropdownLabel = <Icon.Bars />;
        const buttonProps = {
            color: Style.Palette.Secondary,
            outline: true,
            size: Style.Size.XS,
        };
        return (
            <Table.Row>
                <td>{ link }</td>
                <td><SchedulePattern pattern={this.props.schedule.Pattern} /></td>
                {this.enabledOnColumn()}
                <td><DateLabel date={this.props.schedule.LastRunTime} /></td>
                <td><EnabledBadge value={this.props.schedule.Enabled} /></td>
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
