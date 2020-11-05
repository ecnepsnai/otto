import * as React from 'react';
import { NavLink } from 'react-router-dom';
import { Icon } from './Icon';
import { StateManager } from '../services/StateManager';
import { Style } from './Style';
import { UserManager } from '../pages/options/OptionsUsers';
import { API } from '../services/API';
import '../../css/nav.scss';

export class Nav extends React.Component<{}, {}> {
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
                        <div className="collapse navbar-collapse justify-content-between" id="navbarNav">
                            <ul className="navbar-nav">
                                <NavItem link="/hosts" icon={<Icon.Desktop />} label="Hosts" />
                                <NavItem link="/groups" icon={<Icon.LayerGroup />} label="Groups" />
                                <NavItem link="/scripts" icon={<Icon.Scroll />} label="Scripts" />
                                <NavItem link="/schedules" icon={<Icon.Calendar />} label="Schedules" />
                                <NavItem link="/options" icon={<Icon.Cog />} label="Options" />
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
                { StateManager.Current().Warnings.map((warn, idx) => {
                    return ( <Warning warning={warn} key={idx} /> );
                }) }
            </header>
        );
    }
}

interface NavItemProps {
    link: string;
    icon: JSX.Element;
    label: string;
}
class NavItem extends React.Component<NavItemProps, {}> {
    constructor(props: NavItemProps) {
        super(props);
    }

    render(): JSX.Element {
        return (
            <NavLink to={this.props.link} className="nav-link" activeClassName="active">
                <li className="nav-item" data-target=".navbar-collapse.show" data-toggle="collapse">
                    <Icon.Label icon={this.props.icon} label={this.props.label} />
                </li>
            </NavLink>
        );
    }
}

interface WarningProps {
    warning: string;
}
class Warning extends React.Component<WarningProps, {}> {
    render(): JSX.Element {
        let title = '';
        let body = '';

        if (this.props.warning === 'default_user_password') {
            title = 'Default Password';
            body = 'You are using the default username and password. You should change your password immediately using the user menu in the top-right.';
        }

        return (
            <div className="warning">
                <Icon.ExclamationTriangle />
                <strong className="ml-1">
                    {title}
                </strong>
                <span className="ml-1">
                    {body}
                </span>
            </div>
        );
    }
}
