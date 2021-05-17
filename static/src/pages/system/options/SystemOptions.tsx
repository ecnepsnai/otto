import * as React from 'react';
import { Page } from '../../../components/Page';
import { Options } from '../../../types/Options';
import { StateManager } from '../../../services/StateManager';
import { Form } from '../../../components/Form';
import { OptionsGeneral } from './OptionsGeneral';
import { OptionsAuthentication } from './OptionsAuthentication';
import { OptionsNetwork } from './OptionsNetwork';
import { Notification } from '../../../components/Notification';
import { OptionsSecurity } from './OptionsSecurity';

export const SystemOptions: React.FC = () => {
    const [loading, setLoading] = React.useState(false);
    const [options, setOptions] = React.useState(StateManager.Current().Options);

    const changeGeneral = (value: Options.General) => {
        setOptions(options => {
            options.General = value;
            return { ...options };
        });
    };

    const changeAuthentication = (value: Options.Authentication) => {
        setOptions(options => {
            options.Authentication = value;
            return { ...options };
        });
    };

    const changeNetwork = (value: Options.Network) => {
        setOptions(options => {
            options.Network = value;
            return { ...options };
        });
    };

    const changeSecurity = (value: Options.Security) => {
        setOptions(options => {
            options.Security = value;
            return { ...options };
        });
    };

    const onSubmit = () => {
        setLoading(true);
        return Options.Options.Save(options).then(options => {
            Notification.success('Options Saved');
            setOptions(options);
            setLoading(false);
        });
    };

    return (
        <Page title="Options">
            <Form className="cards" showSaveButton onSubmit={onSubmit} loading={loading}>
                <OptionsGeneral defaultValue={options.General} onUpdate={changeGeneral} />
                <OptionsAuthentication defaultValue={options.Authentication} onUpdate={changeAuthentication} />
                <OptionsNetwork defaultValue={options.Network} onUpdate={changeNetwork} />
                <OptionsSecurity defaultValue={options.Security} onUpdate={changeSecurity} />
                <div className="mb-2"></div>
            </Form>
        </Page>
    );
};
