describe('Input', () => {
  it('Textarea input', () => {
    cy.visit('/input')
    cy.get('textarea').type('testarea: 1')
    cy.get('textarea').blur()
    cy.get('textarea').type('2')
    cy.get('textarea').type('3')
    cy.get('textarea').blur()
    cy.contains('Value: testarea: 123')
  })

  it('Textbox input', () => {
    cy.visit('/input')
    cy.get('input[id=textbox_component_Textbox]').type('abc')
    cy.get('input[id=textbox_component_Textbox]').blur()
    cy.get('input[id=textbox_component_Textbox]').type('abc')
    cy.get('input[id=textbox_component_Textbox]').blur()
    cy.contains('Value: abcabc')
  })

  it('Fileupload input', () => {
    cy.visit('/input')
    cy.get('input[type=file]').selectFile('cypress/fixtures/example.json', {
      action: 'drag-drop',
      force: true,
    })
    cy.contains('example.json').should('not.exist')

    cy.get('input[type=file]').selectFile('cypress/fixtures/example.png', {
      action: 'drag-drop',
      force: true,
    })
    cy.contains('example.png').should('exist')
  })

  it('Checkbox', () => {
    cy.visit('/input')
    cy.get('input[type=checkbox]').click()
    cy.contains('Value: true').should('exist')
  })

  it('Button click', () => {
    cy.visit('/input')
    cy.contains('button').click()
    cy.contains('Value: true').should('exist')
    cy.contains('Rerun').click()
    cy.get('div[id=column_component_show_button]').contains('Value: false').should('exist')
  })

  it('Select', () => {
    cy.visit('/input')
    cy.get('select').select(['Value1'])
    cy.get('select').blur()
    cy.contains('Value: Value1').should('exist')

    cy.get('select').select(['Value2'])
    cy.get('select').blur()
    cy.contains('Value: Value2').should('exist')
  })

  it('Radio', () => {
    cy.visit('/input')
    cy.contains('Value3').click()
    cy.contains('Value: Value3').should('exist')

    cy.contains('Value4').click()
    cy.contains('Value: Value4').should('exist')
  })

  it('Datepicker', () => {
    cy.visit('/input')
    cy.get('input[type=date]').type('2000-01-01')
    cy.get('input[type=date]').blur()
    cy.contains('2000-01-01').should('exist')

    cy.get('input[type=date]').type('2002-02-02')
    cy.get('input[type=date]').blur()
    cy.contains('2002-02-02').should('exist')
  })

  it('Timepicker', () => {
    cy.visit('/input')
    cy.get('input[type=time]').type('20:34')
    cy.get('input[type=time]').blur()
    cy.contains('20:34').should('exist')

    cy.get('input[type=time]').type('11:34')
    cy.get('input[type=time]').blur()
    cy.contains('11:34').should('exist')
  })

  it('Datetimepicker', () => {
    cy.visit('/input')
    cy.get('input[type=datetime-local]').type('2000-01-01T20:34')
    cy.get('input[type=datetime-local]').blur()
    cy.contains('2000-01-01 20:34').should('exist')

    cy.get('input[type=datetime-local]').type('2002-01-02T11:34')
    cy.get('input[type=datetime-local]').blur()
    cy.contains('2002-01-02 11:34').should('exist')
  })
})