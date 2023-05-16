import { Checkbox } from './Checkbox';
import { FileChooser } from './FileChooser';
import { RunAsInput } from './RunAsInput';
import { Number } from './Number';
import { Password } from './Password';
import { Radio } from './Radio';
import { Select } from './Select';
import { Text } from './Text';
import { Textarea } from './Textarea';
import { RunLevel } from './RunLevel';

export interface InputProps {
    thin?: boolean;
}

export class Input {
    /**
     * A checkbox and a label.
     */
    public static Checkbox = Checkbox;
    /**
     * A file chooser
     */
    public static FileChooser = FileChooser;
    /**
     * An input node for number type inputs. Suitable for integer and floating points, but not suitable
     * for hexadecimal or scientific-notation values.
     */
    public static Number = Number;
    /**
     * A group of radio buttons with associated labels and values.
     */
    public static Radio = Radio;
    /**
     * A dropdown menu with associated labels and values.
     */
    public static Select = Select;
    /**
     * A text input box. Suitable for plain-text or passwords.
     */
    public static Text = Text;
    /**
     * A multi-line text box.
     */
    public static Textarea = Textarea;
    /**
     * An input for specifying a UID and GID pair
     */
    public static RunAsInput = RunAsInput;
    /**
     * An input for a password with a generate random button
     */
    public static Password = Password;
    /**
     * An input for script run levels
     */
    public static RunLevel = RunLevel;
}
