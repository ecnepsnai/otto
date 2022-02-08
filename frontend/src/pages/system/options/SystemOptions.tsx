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
import { Tabs } from '../../../components/Tabs';
import { Icon } from '../../../components/Icon';

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
            <Form showSaveButton onSubmit={onSubmit} loading={loading}>
                <Tabs.Tabs>
                    <Tabs.Tab icon={<Icon.Wrench />} title="General">
                        <OptionsGeneral defaultValue={options.General} onUpdate={changeGeneral} />
                    </Tabs.Tab>
                    <Tabs.Tab icon={<Icon.Key />} title="Authentication">
                        <OptionsAuthentication defaultValue={options.Authentication} onUpdate={changeAuthentication} />
                    </Tabs.Tab>
                    <Tabs.Tab icon={<Icon.NetworkWired />} title="Network">
                        <OptionsNetwork defaultValue={options.Network} onUpdate={changeNetwork} />
                    </Tabs.Tab>
                    <Tabs.Tab icon={<Icon.Shield />} title="Security">
                        <OptionsSecurity defaultValue={options.Security} onUpdate={changeSecurity} />
                    </Tabs.Tab>
                </Tabs.Tabs>
            </Form>
        </Page>
    );
};
