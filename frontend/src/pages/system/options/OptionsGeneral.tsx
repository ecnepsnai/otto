import * as React from 'react';
import { Input } from '../../../components/input/Input';
import { Options } from '../../../types/Options';
import { EnvironmentVariableEdit } from '../../../components/EnvironmentVariableEdit';
import { Variable } from '../../../types/Variable';
import { Alert } from '../../../components/Alert';

interface OptionsGeneralProps {
    defaultValue: Options.General;
    onUpdate: (value: Options.General) => (void);
}
export const OptionsGeneral: React.FC<OptionsGeneralProps> = (props: OptionsGeneralProps) => {
    const [value, setValue] = React.useState(props.defaultValue);

    React.useEffect(() => {
        props.onUpdate(value);
    }, [value]);

    const originWarning = () => {
        let origin = location.origin;
        if (!origin.endsWith('/')) {
            origin = origin + '/';
        }
        let serverURL = props.defaultValue.ServerURL;
        if (!serverURL.endsWith('/')) {
            serverURL = serverURL + '/';
        }

        if (origin === serverURL) {
            return null;
        }

        return (<Alert.Danger>
            The configured server URL is different than the URL you are using to access this page.
        </Alert.Danger>);
    };

    const changeServerURL = (ServerURL: string) => {
        setValue(value => {
            value.ServerURL = ServerURL;
            return { ...value };
        });
    };

    const changeGlobalEnvironment = (GlobalEnvironment: Variable[]) => {
        setValue(value => {
            value.GlobalEnvironment = GlobalEnvironment;
            return { ...value };
        });
    };

    return (
        <div>
            <Input.Text
                type="text"
                label="Otto Server URL"
                placeholder="https://otto.example.com/"
                helpText="The absolute URL (Including protocol) where this otto server is accessed from"
                defaultValue={value.ServerURL}
                onChange={changeServerURL} />
            {originWarning()}
            <label className="form-label">Global Environment Variables</label>
            <div>
                <EnvironmentVariableEdit
                    variables={value.GlobalEnvironment || []}
                    onChange={changeGlobalEnvironment} />
            </div>
        </div>
    );
};
