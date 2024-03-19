import * as React from 'react';
import { Style } from './Style';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { IconProp } from '@fortawesome/fontawesome-svg-core';
import {
    faArrowDown,
    faArrowLeft,
    faArrowRight,
    faArrowUp,
    faBars,
    faBook,
    faCalendarAlt,
    faCaretDown,
    faCaretUp,
    faCheckCircle,
    faChevronRight,
    faClipboard,
    faCog,
    faDesktop,
    faDownload,
    faEdit,
    faExclamationCircle,
    faExclamationTriangle,
    faExternalLinkAlt,
    faEye,
    faInfoCircle,
    faKey,
    faLayerGroup,
    faLevelDownAlt,
    faListAlt,
    faLock,
    faMagic,
    faMagnifyingGlass,
    faMinus,
    faNetworkWired,
    faPaperclip,
    faPlayCircle,
    faPlus,
    faPuzzlePiece,
    faQuestionCircle,
    faRandom,
    faScroll,
    faShieldAlt,
    faSignOutAlt,
    faSpinner,
    faStarOfLife,
    faTerminal,
    faTimesCircle,
    faTrashAlt,
    faUndo,
    faUnlock,
    faUser,
    faUserEdit,
    faUsers,
    faWrench,
} from '@fortawesome/free-solid-svg-icons';
import '../../css/icon.scss';

export namespace Icon {
    interface IconProps {
        pulse?: boolean;
        spin?: boolean;
        color?: Style.Palette;
    }

    interface EIconProps {
        icon: IconProp;
        options: IconProps;
    }

    export const EIcon: React.FC<EIconProps> = (props: EIconProps) => {
        let className = '';
        if (props.options.color) {
            className = 'text-' + props.options.color.toString();
        }
        return (<FontAwesomeIcon icon={props.icon} pulse={props.options.pulse} spin={props.options.spin} className={className} />);
    };

    interface LabelProps { icon: JSX.Element; spin?: boolean; label: string | number; }
    export const Label: React.FC<LabelProps> = (props: LabelProps) => {
        return (
            <span>
                {props.icon}
                <span className="ms-1">{props.label}</span>
            </span>
        );
    };

    export const ArrowDown: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faArrowDown, options: props });
    export const ArrowLeft: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faArrowLeft, options: props });
    export const ArrowRight: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faArrowRight, options: props });
    export const ArrowUp: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faArrowUp, options: props });
    export const Bars: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faBars, options: props });
    export const Book: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faBook, options: props });
    export const Calendar: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faCalendarAlt, options: props });
    export const CaretDown: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faCaretDown, options: props });
    export const CaretUp: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faCaretUp, options: props });
    export const CheckCircle: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faCheckCircle, options: props });
    export const ChevronRight: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faChevronRight, options: props });
    export const Clipboard: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faClipboard, options: props });
    export const Cog: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faCog, options: props });
    export const Delete: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faTrashAlt, options: props });
    export const Desktop: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faDesktop, options: props });
    export const Download: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faDownload, options: props });
    export const Edit: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faEdit, options: props });
    export const ExclamationCircle: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faExclamationCircle, options: props });
    export const ExclamationTriangle: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faExclamationTriangle, options: props });
    export const ExternalLinkAlt: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faExternalLinkAlt, options: props });
    export const Eye: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faEye, options: props });
    export const InfoCircle: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faInfoCircle, options: props });
    export const Key: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faKey, options: props });
    export const LayerGroup: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faLayerGroup, options: props });
    export const LevelDown: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faLevelDownAlt, options: props });
    export const List: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faListAlt, options: props });
    export const Lock: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faLock, options: props });
    export const Magic: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faMagic, options: props });
    export const MagnifyingGlass: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faMagnifyingGlass, options: props });
    export const Minus: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faMinus, options: props });
    export const NetworkWired: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faNetworkWired, options: props });
    export const Paperclip: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faPaperclip, options: props });
    export const PlayCircle: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faPlayCircle, options: props });
    export const Plus: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faPlus, options: props });
    export const PuzzlePiece: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faPuzzlePiece, options: props });
    export const QuestionCircle: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faQuestionCircle, options: props });
    export const Random: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faRandom, options: props });
    export const Scroll: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faScroll, options: props });
    export const Shield: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faShieldAlt, options: props });
    export const SignOut: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faSignOutAlt, options: props });
    export const Spinner: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faSpinner, options: props });
    export const StarOfLife: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faStarOfLife, options: props });
    export const Terminal: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faTerminal, options: props });
    export const TimesCircle: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faTimesCircle, options: props });
    export const Undo: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faUndo, options: props });
    export const Unlock: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faUnlock, options: props });
    export const User: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faUser, options: props });
    export const UserEdit: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faUserEdit, options: props });
    export const Users: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faUsers, options: props });
    export const Wrench: React.FC<IconProps> = (props: IconProps) => EIcon({ icon: faWrench, options: props });

    // Special icons
    export const Descendant: React.FC = () => {
        return (<span className="descendant-icon"><FontAwesomeIcon icon={faLevelDownAlt} flip="horizontal" transform={{ rotate: 90 }} /></span>);
    };
}
