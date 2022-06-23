import * as React from 'react';
import * as Bootstrap from 'bootstrap';
import { useDebounce } from '@react-hook/debounce';
import { API } from '../services/API';
import { Nothing } from './Nothing';
import { useNavigate } from 'react-router-dom';
import '../../css/system-search.scss';

interface SystemSearchResult {
    Type: 'Host' | 'Group' | 'Script' | 'Schedule' | 'User';
    Label: string;
    URL: string;
}

export const SystemSearch: React.FC = () => {
    const navigate = useNavigate();
    const SEARCH_GROW_WIDTH = 100;
    const SEARCH_ANIMATION_DELAY = 160;

    const [Query, setQuery] = useDebounce('', 300);
    const [Results, setResults] = React.useState<SystemSearchResult[]>();

    const doSearch = () => {
        if (Query == '') {
            setResults(undefined);
            return;
        }

        API.POST('/api/search/system', { Query: Query }).then(data => {
            setResults(data as SystemSearchResult[]);
        });
    };

    React.useEffect(() => {
        if (document.body.offsetWidth > 991) {
            const input = document.getElementById('system-search-input');
            input.style.width = input.offsetWidth + 'px';
        }
    }, []);

    React.useEffect(() => {
        doSearch();
    }, [Query]);

    const onChangeInput = (event: React.ChangeEvent<HTMLInputElement>) => {
        const target = event.target;
        setQuery(target.value);
    };

    const onFocusInput = (event: React.FocusEvent<HTMLInputElement>) => {
        const target = event.target;
        if (document.body.offsetWidth > 991) {
            target.style.width = target.offsetWidth + SEARCH_GROW_WIDTH + 'px';
        }

        if (Query != '') {
            setTimeout(doSearch, SEARCH_ANIMATION_DELAY);
        }
    };

    const onBlurInput = (event: React.FocusEvent<HTMLInputElement>) => {
        setTimeout(() => {
            const target = event.target;
            if (document.body.offsetWidth > 991) {
                target.style.width = target.offsetWidth - SEARCH_GROW_WIDTH + 'px';
            }

            setResults(undefined);
        }, SEARCH_ANIMATION_DELAY);
    };

    const resultClick = (href: string) => {
        new Bootstrap.Collapse(document.getElementById('navbarNav')).hide();
        navigate(href);
        setTimeout(() => {
            setQuery('');
            (document.getElementById('system-search-input') as HTMLInputElement).value = '';
        }, SEARCH_ANIMATION_DELAY);
    };

    return (
        <form className="d-flex">
            <input
                id="system-search-input"
                className="form-control me-2"
                type="search"
                placeholder="Search"
                aria-label="Search"
                defaultValue={Query}
                onChange={onChangeInput}
                onFocus={onFocusInput}
                onBlur={onBlurInput} />
            <ResultsList results={Results} onClick={resultClick} />
        </form>
    );
};

interface ResultsListProps {
    results: SystemSearchResult[];
    onClick: (href: string) => void;
}

const ResultsList: React.FC<ResultsListProps> = (props: ResultsListProps) => {
    const MAX_VISIBLE_RESULTS = 10;

    if (!props.results) {
        return null;
    }

    const input = document.getElementById('system-search-input');
    if (!input) {
        return null;
    }

    const style = {
        top: input.offsetTop + input.offsetHeight + 5,
        left: input.offsetLeft,
        width: input.offsetWidth,
    };

    const linkClick = (href: string) => {
        return () => {
            props.onClick(href);
        };
    };

    const moreResults = () => {
        if (props.results.length < MAX_VISIBLE_RESULTS) {
            return null;
        }

        return (<span className="dropdown-item disabled">Truncating results...</span>);
    };

    const noResults = () => {
        if (props.results.length == 0) {
            return (<span className="dropdown-item disabled"><Nothing /></span>);
        }
    };

    return (
        <ul className="dropdown-menu show" style={style}>
            {noResults()}
            {props.results.map((result, idx) => {
                if (idx > MAX_VISIBLE_RESULTS) {
                    return null;
                }

                return (
                    <li key={idx}>
                        <span className="dropdown-item" onClick={linkClick(result.URL)}>
                            <span className="badge bg-secondary me-1">{result.Type}</span>
                            <span>{result.Label}</span>
                        </span>
                    </li>
                );
            })}
            {moreResults()}
        </ul>
    );
};
