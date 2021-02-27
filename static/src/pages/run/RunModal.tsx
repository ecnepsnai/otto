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
interface RunModalState {
    stage: RunStage;
    selectedHostIDs?: string[];
    finishedHosts: string[];
}
export class RunModal extends React.Component<RunModalProps, RunModalState> {
    constructor(props: RunModalProps) {
        super(props);
        this.state = {
            stage: props.hostIDs ? RunStage.Running : RunStage.Setup,
            selectedHostIDs: props.hostIDs,
            finishedHosts: [],
        };
    }

    private onSelectHostIDs = (hostIDs: string[]) => {
        this.setState({
            selectedHostIDs: hostIDs
        });
    }

    private checkFinished = () => {
        if (this.state.finishedHosts.length <= this.state.selectedHostIDs.length) {
            this.setState({
                stage: RunStage.Finished
            });
        }
    }

    private scriptFinished = (hostID: string) => {
        return () => {
            this.setState(state => {
                state.finishedHosts.push(hostID);
                return state;
            }, () => {
                this.checkFinished();
            });
        };
    }

    private setup = () => {
        return (
            <RunSetup scriptID={this.props.scriptID} onSelectedHosts={this.onSelectHostIDs}/>
        );
    }

    private running = () => {
        return (
            <div className="cards">
                { this.state.selectedHostIDs.map(hostID => {
                    return (
                        <RunScript scriptID={this.props.scriptID} hostID={hostID} key={hostID} onFinished={this.scriptFinished(hostID)}/>
                    );
                })}
            </div>
        );
    }

    private buttons = (): ModalButton[] => {
        if (this.state.stage == RunStage.Running) {
            return [];
        }

        if (this.state.stage == RunStage.Finished) {
            return [
                {
                    label: 'Close',
                }
            ];
        }

        return [
            {
                label: 'Start',
                onClick: () => {
                    this.setState({
                        stage: RunStage.Running,
                    });
                },
                dontDismiss: true,
            }
        ];
    }

    private content = () => {
        if (this.state.stage == RunStage.Setup) {
            return this.setup();
        } else if (this.state.stage == RunStage.Running || this.state.stage == RunStage.Finished) {
            return this.running();
        }
        return null;
    }

    render(): JSX.Element {
        return (
            <Modal title="Run Script" size={Style.Size.L} buttons={this.buttons()} static={this.state.stage === RunStage.Running}>
                { this.content() }
            </Modal>
        );
    }
}
