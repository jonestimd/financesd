import {configure, ShallowWrapper} from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import 'jest-enzyme';
import {clearTimers} from 'mobx-react-lite';

configure({adapter: new Adapter()});

expect.extend({
    toHaveProps: (received: ShallowWrapper, actual: Record<string, unknown>) => {
        const diff = Object.keys(actual).filter((key) => received.prop(key) !== actual[key]);
        if (diff.length === 0) return {
            pass: true,
            message: () => 'expected prop values not to match',
        };
        return {
            pass: false,
            message: () => `expected ${received.name()} to have props\n${diff.map((key) => {
                // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
                return `\t${key}: \x1b[92m${received.prop(key)}\x1b[0m -> \x1b[91m${actual[key]}\x1b[0m`;
            }).join('\n')}`,
        };
    },
});

global.navigator = {
    language: 'en-US',
} as Navigator;

global.fetch = jest.fn();

afterAll(() => clearTimers());
