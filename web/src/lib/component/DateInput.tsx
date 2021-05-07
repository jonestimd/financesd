import React from 'react';
import {Icon, IconButton, TextField, TextFieldProps} from '@material-ui/core';
import classNames from 'classnames';

type IProps = Omit<TextFieldProps, 'value'> & {
    initialValue: string;
    onDateChange: (date: Date, value: string) => void;
    onDateError?: (value: string) => void;
};

const format = /^\d{4}-(0?[1-9]|1[0-9])-(0?[1-9]|[1-2][0-9]|3[01])$/;
const normalize = (value: string) => {
    const parts = value.split('-');
    return `${parts[0].padStart(4, '0')}-${parts[1].padStart(2, '0')}-${parts[2].padStart(2, '0')}`;
};

const DateInput = React.forwardRef<HTMLDivElement, IProps>(({InputProps, initialValue, onDateChange, onDateError, ...props}, ref) => {
    const [date, setDate] = React.useState(new Date(initialValue));
    const [value, setValue] = React.useState(initialValue);
    const [valid, setValid] = React.useState(format.test(value));
    const inputProps: TextFieldProps['InputProps'] = {
        ...InputProps,
        autoFocus: true,
        endAdornment: <IconButton size='small'><Icon>today</Icon></IconButton>,
    };
    const onChange: React.ChangeEventHandler<HTMLInputElement> = (event) => {
        const value = event.currentTarget.value.replace(/[^0-9-]/g, '');
        setValue(value);
        if (format.test(value)) {
            setValid(true);
            const newDate = new Date(normalize(value));
            if (newDate.toISOString() !== date.toISOString()) {
                setDate(newDate);
                onDateChange(newDate, newDate.toISOString().slice(0, 10));
            }
        }
        else {
            setValid(false);
            onDateError?.(value);
        }
    };
    return <TextField ref={ref} {...props} value={value} className={classNames(props.className, 'date', {error: !valid})} InputProps={inputProps}
        onChange={onChange} />;
});

export default DateInput;
