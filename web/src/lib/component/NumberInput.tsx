import React from 'react';
import {TextField, TextFieldProps} from '@material-ui/core';
import classNames from 'classnames';

interface IProps {
    value: string;
    precision: number
    onChange: (value: string) => void;
    startAdornment?: JSX.Element;
}

const isValid = (text: string) => !isNaN(parseFloat(text));

const NumberInput: React.FC<IProps & Partial<Omit<TextFieldProps, 'onChange'>>> = ({value, precision, onChange, startAdornment, ...inputProps}) => {
    const [text, setText] = React.useState(value);
    const [valid, setValid] = React.useState(isValid(value));
    const pattern = React.useMemo(() => new RegExp(`^-?\\d*(\\.\\d{0,${precision}})?$`), [precision]);
    return (
        <TextField {...inputProps} value={text} required InputProps={{autoFocus: true, startAdornment}} className={classNames({error: !valid})}
            onChange={({currentTarget}) => {
                if (pattern.test(currentTarget.value)) {
                    setText(currentTarget.value);
                    setValid(isValid(currentTarget.value));
                    onChange(currentTarget.value);
                }
            }} />
    );
};

export default NumberInput;
