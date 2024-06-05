import * as React from 'react';
import { Input } from './Input';
import { ScriptRunLevel } from '../../types/gengo_enum';

interface RunLevelProps {
    defaultValue: ScriptRunLevel;
    onChange: (runLevel: ScriptRunLevel) => (void);
}
export const RunLevel: React.FC<RunLevelProps> = (props: RunLevelProps) => {
    const [runLevel, setRunLevel] = React.useState<string>(props.defaultValue + '');

    const changeRunLevel = (newLevel: string) => {
        setRunLevel(newLevel);
        props.onChange(parseInt(newLevel) as ScriptRunLevel);
    };

    const choices = [
        {
            label: 'Read-Only',
            value: ScriptRunLevel.ReadOnly + '',
        },
        {
            label: 'Read-Write',
            value: ScriptRunLevel.ReadWrite + '',
        },
    ];

    return (
        <React.Fragment>
            <Input.Radio label='Scope' buttons choices={choices} defaultValue={runLevel} onChange={changeRunLevel} />
        </React.Fragment>
    );
};