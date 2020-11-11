import * as React from 'react';
import { Rand } from '../services/Rand';
import { Icon } from './Icon';

export interface RandomPSKProps {
    newPSK: (psk: string) => (void);
}
export class RandomPSK extends React.Component<RandomPSKProps, {}> {
    private randomPSK = (event: React.MouseEvent<HTMLAnchorElement>) => {
        event.preventDefault();
        this.props.newPSK(Rand.PSK());
    }

    render(): JSX.Element {
        return (
            <div className="mb-3">
                <a href="#" onClick={this.randomPSK} className="mb-3">
                    <Icon.Label icon={<Icon.Random />} label="Generate Random PSK" />
                </a>
            </div>
        );
    }
}
