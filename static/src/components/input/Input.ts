import { Checkbox } from './Checkbox';
import { FileChooser } from './FileChooser';
import { IDInput } from './IDInput';
import { Number } from './Number';
import { Radio } from './Radio';
import { Select } from './Select';
import { Text } from './Text';
import { Textarea } from './Textarea';

export class Input {
    /**
     * A checkbox and a label.
     */
    public static Checkbox = Checkbox
    /**
     * A file chooser
     */
    public static FileChooser = FileChooser;
    /**
     * An input node for number type inputs. Suitable for integer and floating points, but not suitable
     * for hexadecimal or scientific-notation values.
     */
    public static Number = Number
    /**
     * A group of radio buttons with associated labels and values.
     */
    public static Radio = Radio
    /**
     * A dropdown menu with associated labels and values.
     */
    public static Select = Select
    /**
     * A text input box. Suitable for plain-text or passwords.
     */
    public static Text = Text
    /**
     * A multi-line text box.
     */
    public static Textarea = Textarea
    /**
     * An input for specifying a UID and GID pair
     */
    public static IDInput = IDInput
}
