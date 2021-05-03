import React from 'react';
import {Icon, TextField, TextFieldProps} from '@material-ui/core';

type IProps = TextFieldProps & {
    icon: string;
};

const IconInput = React.forwardRef<HTMLDivElement, IProps>(({InputProps, icon, ...props}, ref) => {
    return <TextField ref={ref} {...props} InputProps={{...InputProps, autoFocus: true, startAdornment: <Icon>{icon}</Icon>}} />;
});

export default IconInput;
