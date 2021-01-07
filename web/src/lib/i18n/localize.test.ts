import {translate, defaultKey} from './localize';

describe('localize', () => {
    describe('translate', () => {
        const language = 'es-SP';
        const defaultValue = 'default value';

        beforeEach(() => {
            jest.spyOn(navigator, 'language', 'get').mockReturnValue(language);
        });
        it('returns value from default bundle if no translation', () => {
            const result = translate('entry', {[defaultKey]: {entry: defaultValue}});

            expect(result).toEqual(defaultValue);
        });
        it('returns translation', () => {
            const spanishValue = 'spanish value';

            const result = translate('entry', {
                [defaultKey]: {entry: defaultValue},
                [language]: {entry: spanishValue},
            });

            expect(result).toEqual(spanishValue);
        });
        it('returns key of no value', () => {
            const result = translate('entry', {[defaultKey]: {}});

            expect(result).toEqual('entry');
        });
    });
});
