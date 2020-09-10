import * as React from 'react';
import { Link } from 'react-router-dom';
import { Icon } from './Icon';
import { StateManager } from '../services/StateManager';
import { Style } from './Style';
import { UserManager } from '../pages/options/OptionsUsers';
import '../../css/nav.scss';
import { API } from '../services/API';

interface NavProps {}
interface NavState {
    active: {[id: string]: boolean};
}
export class Nav extends React.Component<NavProps, NavState> {
    constructor(props: NavProps) {
        super(props);
        this.state = {
            active: {
                "hosts": false,
                "groups": false,
                "scripts": false,
                "schedules": false,
                "options": false,
            },
        };
    }
    componentDidMount(): void {
        this.updateNavClass();

        document.body.addEventListener('click', () => {
            setTimeout(() => { this.updateNavClass(); }, 10);
        });
    }
    private updateNavClass = () => {
        this.setState(state => {
            state.active.hosts = location.pathname.indexOf('/hosts') > -1;
            state.active.groups = location.pathname.indexOf('/groups') > -1;
            state.active.scripts = location.pathname.indexOf('/scripts') > -1;
            state.active.schedules = location.pathname.indexOf('/schedules') > -1;
            state.active.options = location.pathname.indexOf('/options') > -1;
            return state;
        });
    }

    private editUserClick = () => {
        UserManager.EditCurrentUser().then(() => { StateManager.Refresh(); });
    }

    private logoutClick = () => {
        API.POST('/api/logout', {}).then(() => {
            location.href = '/login?logged_out';
        }, () => {
            location.href = '/login?logged_out';
        }).catch(() => {
            location.href = '/login?logged_out';
        });
    }

    render(): JSX.Element {
        return (
            <header>
                <nav className="navbar navbar-expand-lg navbar-dark">
                    <div className="container-fluid">
                        <a className="navbar-brand" href="/">
                            <img src="assets/img/full_red.svg" className="brand" height="30" />
                        </a>
                        <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
                            <span className="navbar-toggler-icon"></span>
                        </button>
                        <div className="collapse navbar-collapse d-flex justify-content-between" id="navbarNav">
                            <ul className="navbar-nav">
                                <NavItem link="/hosts" active={this.state.active.hosts} icon={<Icon.Desktop />} label="Hosts" />
                                <NavItem link="/groups" active={this.state.active.groups} icon={<Icon.LayerGroup />} label="Groups" />
                                <NavItem link="/scripts" active={this.state.active.scripts} icon={<Icon.Scroll />} label="Scripts" />
                                <NavItem link="/schedules" active={this.state.active.schedules} icon={<Icon.Calendar />} label="Schedules" />
                                <NavItem link="/options" active={this.state.active.options} icon={<Icon.Cog />} label="Options" />
                            </ul>
                            <ul className="navbar-nav navbar-links">
                                <li className="nav-item dropdown">
                                    <a className="nav-link dropdown-toggle" id="navUserDropdown" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                                        <Icon.Label icon={<Icon.User />} label={StateManager.Current().User.Username} />
                                    </a>
                                    <div className="dropdown-menu dropdown-menu-right navbar-dropdown" aria-labelledby="navUserDropdown">
                                        <a className="dropdown-item" onClick={this.editUserClick}><Icon.Label icon={<Icon.UserEdit />} label="Edit User" /></a>
                                        <a className="dropdown-item" onClick={this.logoutClick}><Icon.Label icon={<Icon.SignOut />} label="Log Out" /></a>
                                        <div className="dropdown-divider"></div>
                                        <h6 className="dropdown-header">Otto {StateManager.Current().Runtime.Version}</h6>
                                        <a className="dropdown-item" href={'https://github.com/ecnepsnai/otto/tree/' + StateManager.Current().Runtime.Version + '/docs'} target="_blank" rel="noreferrer">
                                            <Icon.Label icon={<Icon.InfoCircle color={Style.Palette.Primary} />} label="Documentation" />
                                        </a>
                                        <a className="dropdown-item" href="https://github.com/ecnepsnai/otto/issues/new" target="_blank" rel="noreferrer">
                                            <Icon.Label icon={<Icon.ExclamationCircle color={Style.Palette.Danger} />} label="Report an Issue" />
                                        </a>
                                    </div>
                                </li>
                            </ul>
                        </div>
                    </div>
                </nav>
            </header>
        );
    }
}

interface NavItemProps {
    link: string;
    icon: JSX.Element;
    label: string;
    active: boolean;
}
class NavItem extends React.Component<NavItemProps, {}> {
    render(): JSX.Element {
        let className = 'nav-link';
        if (this.props.active) {
            className += ' active';
        }
        return (
        <li className="nav-item">
            <Link to={this.props.link} className={className} data-target=".navbar-collapse.show" data-toggle="collapse">
                <Icon.Label icon={this.props.icon} label={this.props.label} />
            </Link>
        </li>
        );
    }
}
