import React from 'react';
import {Icon, IconButton, TextField, TextFieldProps} from '@material-ui/core';
import classNames from 'classnames';

const DateInput = React.forwardRef<HTMLDivElement, TextFieldProps>(({InputProps, ...props}, ref) => {
    const inputProps: TextFieldProps['InputProps'] = {
        ...InputProps,
        autoFocus: true,
        endAdornment: <IconButton size='small'><Icon>today</Icon></IconButton>,
    };
    return <TextField ref={ref} {...props} className={classNames(props.className, 'date')} InputProps={inputProps} />;
});

export default DateInput;
