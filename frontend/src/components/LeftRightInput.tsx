import * as React from 'react';
import { Rand } from '../services/Rand';
import { Button } from './Button';
import { Style } from './Style';
import { Icon } from './Icon';
import '../../css/leftright.scss';

interface LeftRightInputChoice {
    label: string
    value: string
}

interface LeftRightInputProps {
    leftLabel: string | JSX.Element;
    rightLabel: string | JSX.Element;
    choices: LeftRightInputChoice[]
    selected: string[]
    onChange: (selected: string[]) => (void)
}
export const LeftRightInput: React.FC<LeftRightInputProps> = (props: LeftRightInputProps) => {
    const [Selected, SetSelected] = React.useState<string[]>(props.selected);
    const [LeftSelected, SetLeftSelected] = React.useState<string[]>([]);
    const [RightSelected, SetRightSelected] = React.useState<string[]>([]);
    const leftId = Rand.ID();
    const rightId = Rand.ID();

    React.useEffect(() => {
        if (Selected === undefined) {
            return;
        }

        props.onChange(Selected);
    }, [Selected]);

    const getLeftChoices = () => {
        return props.choices.filter(c => {
            return !Selected.includes(c.value);
        });
    };

    const getRightChoices = () => {
        const selections: LeftRightInputChoice[] = [];

        Selected.forEach(s => {
            props.choices.forEach(c => {
                if (s === c.value) {
                    selections.push(c);
                }
            });
        });

        return selections;
    };

    const onLeftSelect = (event: React.FormEvent<HTMLSelectElement>) => {
        const target = event.target as HTMLSelectElement;
        
        const selectedValues: string[] = [];
        for (let i = 0; i < target.selectedOptions.length; i++) {
            selectedValues.push(target.selectedOptions[i].value);
        }

        SetLeftSelected(selectedValues);
    };

    const onRightSelect = (event: React.FormEvent<HTMLSelectElement>) => {
        const target = event.target as HTMLSelectElement;
        
        const selectedValues: string[] = [];
        for (let i = 0; i < target.selectedOptions.length; i++) {
            selectedValues.push(target.selectedOptions[i].value);
        }

        SetRightSelected(selectedValues);
    };

    const moveLeftToRight = () => {
        SetSelected(sel => {
            const newSel = sel.concat(LeftSelected);
            return [...newSel];
        });
        SetLeftSelected([]);
        SetRightSelected([]);
    };

    const moveRightToLeft = () => {
        SetSelected(sel => {
            const newSel = sel.filter(s => {
                return !RightSelected.includes(s);
            });
            return [...newSel];
        });
        SetLeftSelected([]);
        SetRightSelected([]);
    };

    const moveUp = () => {
        SetSelected(sel => {
            for (let i = 0; i < RightSelected.length; i++) {
                let index = sel.indexOf(RightSelected[i]);
                sel.splice(index, 1);
                if (index > 0) {
                    index--;
                }
                sel.splice(index, 0, RightSelected[i]);
            }
            
            return [...sel];
        });
    };

    const moveDown = () => {
        SetSelected(sel => {
            for (let i = RightSelected.length-1; i >= 0; i--) {
                let index = sel.indexOf(RightSelected[i]);
                sel.splice(index, 1);
                if (index < sel.length) {
                    index++;
                }
                sel.splice(index, 0, RightSelected[i]);
            }
            
            return [...sel];
        });
    };

    return (
        <div className="container text-center">
            <div className="row">
                <div className="col-5 text-md-start">
                    { props.leftLabel }
                    <select id={leftId} className="form-control" multiple value={LeftSelected} onChange={onLeftSelect}>
                        {
                            getLeftChoices().map(choice => {
                                return (<option key={choice.value} value={choice.value}>{choice.label}</option>);
                            })
                        }
                    </select>
                </div>
                <div className="col lr-buttons mt-4">
                    <Button onClick={moveLeftToRight} size={Style.Size.S} color={Style.Palette.Primary} outline disabled={LeftSelected.length == 0}><Icon.ArrowRight /></Button>
                    <Button onClick={moveRightToLeft} size={Style.Size.S} color={Style.Palette.Primary} outline disabled={LeftSelected.length == 0}><Icon.ArrowLeft /></Button>
                </div>
                <div className="col-5 text-md-start">
                { props.rightLabel }
                    <select id={rightId} className="form-control" multiple value={RightSelected} onChange={onRightSelect}>
                        {
                            getRightChoices().map(choice => {
                                return (<option key={choice.value} value={choice.value}>{choice.label}</option>);
                            })
                        }
                    </select>
                </div>
                <div className="col lr-buttons mt-4">
                    <Button onClick={moveUp} size={Style.Size.S} color={Style.Palette.Primary} outline disabled={RightSelected.length == 0}><Icon.ArrowUp /></Button>
                    <Button onClick={moveDown} size={Style.Size.S} color={Style.Palette.Primary} outline disabled={RightSelected.length == 0}><Icon.ArrowDown /></Button>
                </div>
            </div>
        </div>
    );
};
