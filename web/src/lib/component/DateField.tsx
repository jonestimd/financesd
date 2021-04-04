import React from 'react';
import {Icon, IconButton, InputAdornment, TextField, TextFieldProps} from '@material-ui/core';
import classNames from 'classnames';

const DateField = React.forwardRef<HTMLDivElement, TextFieldProps>((props, ref) => {
    const inputProps: TextFieldProps['InputProps'] = {
        ...props.InputProps,
        endAdornment: <InputAdornment position='end'><IconButton size='small'><Icon>today</Icon></IconButton></InputAdornment>,
    };
    return <TextField ref={ref} {...props} className={classNames(props.className, 'date')} InputProps={inputProps} />;
});

export default DateField;
