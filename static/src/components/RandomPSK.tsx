import * as React from 'react';
import { Rand } from '../services/Rand';
import { Icon } from './Icon';

export interface RandomPSKProps {
    newPSK: (psk: string) => (void);
}
export const RandomPSK: React.FC<RandomPSKProps> = (props: RandomPSKProps) => {
    const randomPSK = (event: React.MouseEvent<HTMLAnchorElement>) => {
        event.preventDefault();
        props.newPSK(Rand.PSK());
    };

    return (
        <div className="mb-3">
            <a href="#" onClick={randomPSK} className="mb-3">
                <Icon.Label icon={<Icon.Random />} label="Generate Random PSK" />
            </a>
        </div>
    );
};
