import * as React from 'react';
import { Modal, ModalButton } from '../../components/Modal';
import { RunSetup } from './RunSetup';
import { RunScript } from './RunScript';
import { Style } from '../../components/Style';

enum RunStage {
    Setup,
    Running,
    Finished,
}

interface RunModalProps {
    scriptID: string;
    hostIDs?: string[];
}
export const RunModal: React.FC<RunModalProps> = (props: RunModalProps) => {
    const [stage, setStage] = React.useState(props.hostIDs ? RunStage.Running : RunStage.Setup);
    const [selectedHostIDs, setSelectedHostIDs] = React.useState(props.hostIDs);
    const [finishedHosts, setFinishedHosts] = React.useState<string[]>([]);

    const onSelectHostIDs = (hostIDs: string[]) => {
        setSelectedHostIDs(hostIDs);
    };

    React.useEffect(() => {
        checkFinished();
    }, [finishedHosts]);

    const checkFinished = () => {
        if (selectedHostIDs && finishedHosts.length <= selectedHostIDs.length) {
            setStage(RunStage.Finished);
        }
    };

    const scriptFinished = (hostID: string) => {
        return () => {
            setFinishedHosts(finishedHosts => {
                return [...finishedHosts, hostID];
            });
        };
    };

    const setup = () => {
        return (
            <RunSetup scriptID={props.scriptID} onSelectedHosts={onSelectHostIDs} />
        );
    };

    const running = () => {
        return (
            <div className="cards">
                {selectedHostIDs.map(hostID => {
                    return (
                        <RunScript scriptID={props.scriptID} hostID={hostID} key={hostID} onFinished={scriptFinished(hostID)} />
                    );
                })}
            </div>
        );
    };

    const buttons = (): ModalButton[] => {
        if (stage == RunStage.Running) {
            return [];
        }

        if (stage == RunStage.Finished) {
            return [
                {
                    label: 'Close',
                }
            ];
        }

        return [
            {
                label: 'Cancel',
                color: Style.Palette.Secondary,
            },
            {
                label: 'Start',
                onClick: () => {
                    if (selectedHostIDs && selectedHostIDs.length > 0) {
                        setStage(RunStage.Running);
                    }
                },
                dontDismiss: true,
            }
        ];
    };

    const content = () => {
        if (stage == RunStage.Setup) {
            return setup();
        } else if (stage == RunStage.Running || stage == RunStage.Finished) {
            return running();
        }
        return null;
    };

    return (
        <Modal title="Run Script" size={Style.Size.L} buttons={buttons()} static>
            {content()}
        </Modal>
    );
};
