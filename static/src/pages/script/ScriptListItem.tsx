import * as React from 'react';
import { Script } from '../../types/Script';
import { Link } from 'react-router-dom';
import { MenuItem, MenuLink } from '../../components/Menu';
import { Icon } from '../../components/Icon';
import { Table } from '../../components/Table';
import { GlobalModalFrame } from '../../components/Modal';
import { Rand } from '../../services/Rand';
import { RunModal } from '../run/RunModal';

export interface ScriptListItemProps { script: Script, onReload: () => (void); }
export class ScriptListItem extends React.Component<ScriptListItemProps, {}> {
    private deleteMenuClick = () => {
        this.props.script.DeleteModal().then(confirmed => {
            if (confirmed) {
                this.props.onReload();
            }
        });
    }

    private toggleMenuClick = () => {
        this.props.script.Update({
            Enabled: !this.props.script.Enabled
        }).then(() => {
            this.props.onReload();
        });
    }

    private enableDisableMenu = () => {
        if (this.props.script.Enabled) {
            return ( <MenuItem icon={<Icon.TimesCircle />} onClick={this.toggleMenuClick} label="Disable" /> );
        }
        return ( <MenuItem icon={<Icon.CheckCircle />} onClick={this.toggleMenuClick} label="Enable" /> );
    }

    private executeScriptMenuClick = () => {
        GlobalModalFrame.showModal(<RunModal scriptID={this.props.script.ID} key={Rand.ID()}/>);
    }

    render(): JSX.Element {
        const link = <Link to={'/scripts/script/' + this.props.script.ID}>{ this.props.script.Name }</Link>;

        return (
            <Table.Row disabled={!this.props.script.Enabled}>
                <td>{ link }</td>
                <td>{ this.props.script.Executable }</td>
                <Table.Menu>
                    <MenuItem label="Run Script" icon={<Icon.PlayCircle />} onClick={this.executeScriptMenuClick}/>
                    <MenuLink label="Edit" icon={<Icon.Edit />} to={'/scripts/script/' + this.props.script.ID + '/edit'}/>
                    { this.enableDisableMenu() }
                    <MenuItem label="Delete" icon={<Icon.Delete />} onClick={this.deleteMenuClick}/>
                </Table.Menu>
            </Table.Row>
        );
    }
}
