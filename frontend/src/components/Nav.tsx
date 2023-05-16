import * as React from 'react';
import { API } from '../services/API';
import { Collapse } from 'bootstrap';
import { Icon } from './Icon';
import { Link, NavLink } from 'react-router-dom';
import { StateManager } from '../services/StateManager';
import { Style } from './Style';
import { UserManager } from '../pages/system/users/SystemUsers';
import { SystemSearch } from './SystemSearch';
import { Permissions, UserAction } from '../services/Permissions';
import '../../css/nav.scss';

export const Nav: React.FC = () => {
    return (
        <header>
            <nav className="navbar navbar-expand-lg navbar-dark">
                <div className="container-fluid">
                    <a className="navbar-brand" href="/">
                        <img src="assets/img/full_red.svg" className="brand" height="30" />
                    </a>
                    <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
                        <span className="navbar-toggler-icon"></span>
                    </button>
                    <div className="collapse navbar-collapse justify-content-between" id="navbarNav">
                        <ul className="navbar-nav">
                            <NavItem link="/hosts" icon={<Icon.Desktop />} label="Hosts" />
                            <NavItem link="/groups" icon={<Icon.LayerGroup />} label="Groups" />
                            <NavItem link="/scripts" icon={<Icon.Scroll />} label="Scripts" />
                            <NavItem link="/schedules" icon={<Icon.Calendar />} label="Schedules" />
                            { Permissions.UserCan(UserAction.AccessAuditLog) ? <NavItem link="/events" icon={<Icon.List />} label="Event Log" /> : null }
                        </ul>
                        <div className="d-flex">
                            <SystemSearch />
                            <SystemMenu />
                            <UserMenu />
                        </div>
                    </div>
                </div>
            </nav>
            { StateManager.Current().Warnings.map((warn, idx) => {
                return (<Warning warning={warn} key={idx} />);
            })}
        </header>
    );
};

interface NavItemProps {
    link: string;
    icon: JSX.Element;
    label: string;
}
export const NavItem: React.FC<NavItemProps> = (props: NavItemProps) => {
    const didClick = () => {
        if (document.body.offsetWidth < 992) {
            const bsc = new Collapse(document.getElementById('navbarNav'));
            bsc.hide();
        }
    };

    const className = (p: { isActive: boolean; }) => {
        return 'nav-link ' + (p.isActive ? 'active' : '');
    };

    return (
        <NavLink to={props.link} onClick={didClick} className={className}>
            <li className="nav-item" data-target=".navbar-collapse.show" data-toggle="collapse">
                <Icon.Label icon={props.icon} label={props.label} />
            </li>
        </NavLink>
    );
};

export const SystemMenu: React.FC = () => {
    const optionsMenu = () => {
        if (!Permissions.UserCan(UserAction.ModifySystem)) {
            return (<span className="dropdown-item disabled"><Icon.Label icon={<Icon.Wrench />} label="Options" /></span>);
        }

        return (<Link to="/system/options" className="dropdown-item"><Icon.Label icon={<Icon.Wrench />} label="Options" /></Link>);
    };

    const usersMenu = () => {
        if (!Permissions.UserCan(UserAction.ModifyUsers)) {
            return (<span className="dropdown-item disabled"><Icon.Label icon={<Icon.User />} label="Users" /></span>);
        }

        return (<Link to="/system/users" className="dropdown-item"><Icon.Label icon={<Icon.User />} label="Users" /></Link>);
    };

    const registerMenu = () => {
        if (!Permissions.UserCan(UserAction.ModifyAutoregister)) {
            return (<span className="dropdown-item disabled"><Icon.Label icon={<Icon.Magic />} label="Host Registration" /></span>);
        }

        return (<Link to="/system/register" className="dropdown-item"><Icon.Label icon={<Icon.Magic />} label="Host Registration" /></Link>);
    };

    return (
        <ul className="navbar-nav navbar-links">
            <li className="nav-item dropdown">
                <a className="nav-link dropdown-toggle" id="navSystemDropdown" role="button" data-bs-toggle="dropdown" aria-expanded="false">
                    <Icon.Label icon={<Icon.Cog />} label="System" />
                </a>
                <div className="dropdown-menu dropdown-menu-end" aria-labelledby="navSystemDropdown">
                    { optionsMenu() }
                    { usersMenu() }
                    { registerMenu() }
                    <div className="dropdown-divider"></div>
                    <h6 className="dropdown-header">Otto {StateManager.Current().Runtime.Version}</h6>
                    <a className="dropdown-item" href={'https://github.com/ecnepsnai/otto/tree/' + StateManager.Current().Runtime.Version + '/docs'} target="_blank" rel="noreferrer">
                        <Icon.Label icon={<Icon.InfoCircle color={Style.Palette.Primary} />} label="Documentation" />
                    </a>
                    <a className="dropdown-item" href="https://github.com/ecnepsnai/otto/issues/new/choose" target="_blank" rel="noreferrer">
                        <Icon.Label icon={<Icon.ExclamationCircle color={Style.Palette.Danger} />} label="Report an Issue" />
                    </a>
                </div>
            </li>
        </ul>
    );
};

export const UserMenu: React.FC = () => {
    const editUserClick = () => {
        UserManager.EditCurrentUser().then(() => {
            StateManager.Refresh();
        });
    };

    const logoutClick = () => {
        API.POST('/api/logout', {}).then(() => {
            location.href = '/login?logged_out';
        }, () => {
            location.href = '/login?logged_out';
        }).catch(() => {
            location.href = '/login?logged_out';
        });
    };

    return (
        <ul className="navbar-nav navbar-links">
            <li className="nav-item dropdown">
                <a className="nav-link dropdown-toggle" id="navUserDropdown" role="button" data-bs-toggle="dropdown" aria-expanded="false">
                    <Icon.Label icon={<Icon.User />} label={StateManager.Current().User.Username} />
                </a>
                <div className="dropdown-menu dropdown-menu-end" aria-labelledby="navUserDropdown">
                    <a className="dropdown-item" onClick={editUserClick}><Icon.Label icon={<Icon.UserEdit />} label="Edit User" /></a>
                    <a className="dropdown-item" onClick={logoutClick}><Icon.Label icon={<Icon.SignOut />} label="Log Out" /></a>
                </div>
            </li>
        </ul>
    );
};

interface WarningProps {
    warning: 'default_user_password';
}
export const Warning: React.FC<WarningProps> = (props: WarningProps) => {
    let title = '';
    let body = '';

    if (props.warning === 'default_user_password') {
        title = 'Default Password';
        body = 'You are using the default username and password. You should change your password immediately using the user menu in the top-right.';
    }

    return (
        <div className="warning">
            <Icon.ExclamationTriangle />
            <strong className="ms-1">
                {title}
            </strong>
            <span className="ms-1">
                {body}
            </span>
        </div>
    );
};
