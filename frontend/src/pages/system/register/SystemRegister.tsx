import * as React from 'react';
import { Button } from '../../../components/Button';
import { Card } from '../../../components/Card';
import { Input } from '../../../components/input/Input';
import { Form } from '../../../components/Form';
import { Icon } from '../../../components/Icon';
import { PageLoading } from '../../../components/Loading';
import { GlobalModalFrame, Modal } from '../../../components/Modal';
import { Notification } from '../../../components/Notification';
import { Page } from '../../../components/Page';
import { Rand } from '../../../services/Rand';
import { StateManager } from '../../../services/StateManager';
import { Group, GroupType } from '../../../types/Group';
import { Options } from '../../../types/Options';
import { RegisterRuleType, RegisterRule } from '../../../types/RegisterRule';
import { Style } from '../../../components/Style';
import { Pre } from '../../../components/Pre';
import { RegisterRules } from './RegisterRules';
import '../../../../css/registerrules.scss';

export const SystemRegister: React.FC = () => {
    const [loading, setLoading] = React.useState<boolean>(true);
    const [rules, setRules] = React.useState<RegisterRuleType[]>();
    const [groups, setGroups] = React.useState<GroupType[]>();
    const [options, setOptions] = React.useState<Options.Register>();
    const [keyInputKey, setKeyInputKey] = React.useState(Rand.ID());

    React.useEffect(() => {
        loadData();
    }, []);

    const loadRules = () => {
        return RegisterRule.List().then(rules => {
            setRules(rules);
        });
    };

    const loadGroups = () => {
        return Group.List().then(groups => {
            setGroups(groups);
        });
    };

    const loadOptions = () => {
        return Options.Options.Get().then(o => {
            setOptions(o.Register);
        });
    };

    const loadData = () => {
        Promise.all([loadRules(), loadGroups(), loadOptions()]).then(() => {
            setLoading(false);
        });
    };

    const onSubmit = () => {
        return Options.Options.Get().then(o => {
            o.Register = options;
            Options.Options.Save(o).then(() => {
                Notification.success('Changes Saved');
            });
        });
    };

    const changeEnabled = (Enabled: boolean) => {
        setOptions(options => {
            options.Enabled = Enabled;
            return { ...options };
        });
    };

    const changeKey = (Key: string) => {
        setOptions(options => {
            options.Key = Key;
            return { ...options };
        });
    };

    const randomKey = () => {
        changeKey(Rand.PSK());
        setKeyInputKey(Rand.ID());
    };

    const showCommand = () => {
        const registerKey = options.Key.replace('"', '\\"');
        let url = StateManager.Current().Options.General.ServerURL;
        url = url.substr(0, url.length - 1);
        const command = 'REGISTER_KEY="' + registerKey + '" \\\nREGISTER_HOST="' + url + '" \\\n/opt/otto-agent/agent';

        GlobalModalFrame.showModal(<Modal title="Register Command">
            <p>Copy and paste the following command into any shell to start agent registration</p>
            <Pre>{command}</Pre>
        </Modal>);
    };

    const changeRunScriptsOnRegister = (RunScriptsOnRegister: boolean) => {
        setOptions(options => {
            options.RunScriptsOnRegister = RunScriptsOnRegister;
            return { ...options };
        });
    };

    const addRule = (rule: RegisterRuleType) => {
        RegisterRule.New(rule).then(() => {
            Notification.success('Rule Added');
            loadRules().then(() => {
                setLoading(false);
            });
        });
    };

    const modifyRule = (id: string, rule: RegisterRuleType) => {
        RegisterRule.Save(id, rule).then(() => {
            Notification.success('Rule Modified');
            loadRules().then(() => {
                setLoading(false);
            });
        });
    };

    const deleteRule = (rule: RegisterRuleType) => {
        RegisterRule.DeleteModal(rule).then(confirmed => {
            if (!confirmed) {
                return;
            }

            loadRules().then(() => {
                setLoading(false);
            });
        });
    };

    const changeDefaultGroupID = (DefaultGroupID: string) => {
        setOptions(options => {
            options.DefaultGroupID = DefaultGroupID;
            return { ...options };
        });
    };

    const enabledContent = () => {
        if (!options.Enabled) {
            return null;
        }

        return (<React.Fragment>
            <Input.Text
                type="text"
                label="Register Key"
                helpText="Agents that wish to register with this server must specify this key to authenticate"
                defaultValue={options.Key}
                onChange={changeKey}
                key={keyInputKey}
                required />
            <div className="buttons mb-2">
                <Button color={Style.Palette.Secondary} size={Style.Size.XS} outline onClick={randomKey}><Icon.Label icon={<Icon.Random />} label="Generate Random Key" /></Button>
                <Button color={Style.Palette.Secondary} size={Style.Size.XS} outline onClick={showCommand}><Icon.Label icon={<Icon.Terminal />} label="Show Command Line" /></Button>
            </div>
            <Card.Card className="mb-2">
                <Card.Header>
                    Rules
                </Card.Header>
                <Card.Body>
                    <RegisterRules rules={rules} onAdd={addRule} onChange={modifyRule} onDelete={deleteRule} groups={groups} />
                </Card.Body>
            </Card.Card>
            <Input.Select
                label="Default Group"
                helpText="If none of the above rules match the agent will be added to this group"
                defaultValue={options.DefaultGroupID}
                onChange={changeDefaultGroupID}>
                {groups.map((group, idx) => {
                    return (<option key={idx} value={group.ID}>{group.Name}</option>);
                })}
            </Input.Select>
            <Input.Checkbox
                label="Automatically Run Scripts on Registration"
                helpText="When a host successfully registers all scripts associated with the host will run automatically."
                defaultValue={options.RunScriptsOnRegister}
                onChange={changeRunScriptsOnRegister} />
        </React.Fragment>);
    };

    if (loading) {
        return (<PageLoading />);
    }

    return (
        <Page title="Host Registration">
            <Form showSaveButton={true} onSubmit={onSubmit}>
                <Input.Checkbox
                    label="Allow Hosts to Register Themselves"
                    helpText="If checked hosts can automatically register themselves with this Otto server"
                    defaultValue={options.Enabled}
                    onChange={changeEnabled} />
                {enabledContent()}
            </Form>
        </Page>
    );
};
