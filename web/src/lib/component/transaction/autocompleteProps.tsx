import React from 'react';
import {AutocompleteProps} from '@material-ui/lab/Autocomplete';

const autocompleteProps: Pick<AutocompleteProps<unknown, false, false, false>, 'openOnFocus' | 'onHighlightChange'> = {
    openOnFocus: true,
    onHighlightChange: function(event: React.ChangeEvent<unknown>, _option: unknown, reason: string) {
        if (reason === 'keyboard') event.stopPropagation();
    },
};

export default autocompleteProps;
