import React from 'react';
import autocompleteProps from './autocompleteProps';

describe('autocompleteProps', () => {
    describe('onHighlightChange', () => {
        it('stops event propagation if reason is keyboard', () => {
            const event = {stopPropagation: jest.fn()} as unknown as React.ChangeEvent;

            autocompleteProps.onHighlightChange!(event, undefined, 'keyboard');
        });
        it('does not stop event propagation if reason is not keyboard', () => {
            const event = {stopPropagation: jest.fn()} as unknown as React.ChangeEvent;

            autocompleteProps.onHighlightChange!(event, undefined, 'mouse');
        });
    });
});
