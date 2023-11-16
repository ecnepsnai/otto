import * as React from 'react';
import { Options } from '../../../types/Options';
import { Card } from '../../../components/Card';
import { Loading } from '../../../components/Loading';
import { Button } from '../../../components/Button';
import { Style } from '../../../components/Style';
import { StateManager } from '../../../services/StateManager';
import { Icon } from '../../../components/Icon';

export const OptionsAdvanced: React.FC = () => {
    const [IsLoading, SetIsLoading] = React.useState(false);

    const setVerboseLogging = (enabled: boolean) => {
        return () => {
            SetIsLoading(true);
            Options.Options.SetVerboseLogging(enabled).then(() => {
                SetIsLoading(false);
            });
        };
    };

    if (IsLoading) {
        return (
            <div>
                <Card.Card>
                    <Card.Header>Verbose Logging</Card.Header>
                    <Card.Body>
                        <Loading />
                    </Card.Body>
                </Card.Card>
            </div>
        );
    }

    return (
        <div>
            <Card.Card>
                <Card.Header>Verbose Logging</Card.Header>
                <Card.Body>
                    <p>By default the Otto server only records informational or error messages to the log. By enabling verbose logging, the Otto server will record much more information to the log. This can be used to troubleshoot issues when communicating with hosts or performing operations on the server.</p>
                    <p>Current status: { StateManager.Current().Runtime.Verbose ? (<Icon.Label icon={<Icon.ExclamationTriangle color={Style.Palette.Warning} />} label='Enabled' />) : (<Icon.Label icon={<Icon.CheckCircle color={Style.Palette.Success} />} label='Disabled' />) }</p>
                    { StateManager.Current().Runtime.Verbose ? (<Button color={Style.Palette.Primary} onClick={setVerboseLogging(false)}>Disable</Button>) : (<Button color={Style.Palette.Primary} onClick={setVerboseLogging(true)}>Enable</Button>) }
                </Card.Body>
            </Card.Card>
        </div>
    );
};
