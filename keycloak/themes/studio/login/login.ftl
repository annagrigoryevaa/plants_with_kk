<#import "template.ftl" as layout>

<@layout.registrationLayout displayMessage=false displayInfo=false; section>
  <#if section = "header">
  <#elseif section = "form">
    <div class="app-wrapper">
      <div class="app-content">
        <div class="login-page">
          <main class="login-page__main">
            <div class="login-page__stack">
              <img class="login-page__logo" src="${url.resourcesPath}/img/logo.svg" alt="GigaChat" />

              <h1 class="login-page__title">
                <span class="login-page__title-platform">GigaPlatform </span>AI Studio
              </h1>

              <p class="login-page__subtitle">Войдите, чтобы начать</p>

              <ul class="login-page__features">
                <li class="login-page__feature">
                  <div aria-hidden="true" class="login-page__feature-icon">
                    <svg width="100%" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                      <path fill-rule="evenodd" clip-rule="evenodd" d="M8.39702 12.6857C7.88593 12.8335 7.34735 12.9166 6.79073 12.9254C7.63838 13.5691 8.69563 13.9511 9.84211 13.9511C10.3103 13.9511 10.7643 13.8873 11.1957 13.7675C11.6317 13.6465 12.0329 13.6052 12.387 13.6984L13.7445 14.0556C14.5036 14.2554 15.1964 13.5626 14.9966 12.8035L14.6394 11.446C14.5462 11.0919 14.5875 10.6907 14.7085 10.2547C14.8283 9.82337 14.8921 9.3693 14.8921 8.90113C14.8921 7.03936 13.8846 5.41289 12.3851 4.53719C12.5827 5.03354 12.7172 5.56182 12.7788 6.11214C13.4687 6.83839 13.8921 7.82032 13.8921 8.90113C13.8921 9.27805 13.8408 9.64216 13.745 9.98723C13.6023 10.5011 13.5146 11.1014 13.6723 11.7005L14.0295 13.058C14.0308 13.0629 14.0309 13.067 14.0309 13.067L14.0304 13.0706C14.0296 13.0732 14.0275 13.0775 14.023 13.0821C14.0185 13.0866 14.0141 13.0886 14.0116 13.0894L14.008 13.0899C14.008 13.0899 14.0039 13.0898 13.999 13.0886L12.6415 12.7313C12.0424 12.5736 11.4421 12.6613 10.9282 12.804C10.5831 12.8998 10.219 12.9511 9.84211 12.9511C9.33301 12.9511 8.84585 12.8572 8.39702 12.6857Z" fill="currentColor"></path>
                      <path fill-rule="evenodd" clip-rule="evenodd" d="M6.65023 2.74951C4.41347 2.74951 2.60022 4.56276 2.60022 6.79951C2.60022 7.17642 2.65158 7.54053 2.74737 7.88561C2.89003 8.39949 2.97771 8.99976 2.82004 9.59892L2.46281 10.9564C2.46152 10.9613 2.46145 10.9654 2.46145 10.9654L2.46198 10.969C2.46272 10.9715 2.4648 10.9759 2.46931 10.9804C2.47381 10.9849 2.4782 10.987 2.48077 10.9878L2.48437 10.9883C2.48437 10.9883 2.48847 10.9882 2.49335 10.9869L3.85082 10.6297C4.44998 10.472 5.05025 10.5597 5.56413 10.7024C5.9092 10.7982 6.27331 10.8495 6.65023 10.8495C8.88698 10.8495 10.7002 9.03626 10.7002 6.79951C10.7002 4.56276 8.88698 2.74951 6.65023 2.74951ZM1.60022 6.79951C1.60022 4.01047 3.86119 1.74951 6.65023 1.74951C9.43926 1.74951 11.7002 4.01047 11.7002 6.79951C11.7002 9.58855 9.43926 11.8495 6.65023 11.8495C6.18205 11.8495 5.72799 11.7857 5.29663 11.6659C4.86066 11.5449 4.45947 11.5036 4.10531 11.5968L2.74784 11.954C1.98871 12.1538 1.29596 11.461 1.49573 10.7019L1.85296 9.34442C1.94616 8.99027 1.90484 8.58907 1.78381 8.1531C1.66407 7.72175 1.60022 7.26768 1.60022 6.79951Z" fill="currentColor"></path>
                      <path fill-rule="evenodd" clip-rule="evenodd" d="M6.65022 4.89941C6.92637 4.89941 7.15022 5.12327 7.15022 5.39941V8.19941C7.15022 8.47556 6.92637 8.69941 6.65022 8.69941C6.37408 8.69941 6.15022 8.47556 6.15022 8.19941V5.39941C6.15022 5.12327 6.37408 4.89941 6.65022 4.89941Z" fill="currentColor"></path>
                      <path fill-rule="evenodd" clip-rule="evenodd" d="M9.10013 5.5998C9.37627 5.5998 9.60013 5.82366 9.60013 6.0998V7.4998C9.60013 7.77595 9.37627 7.9998 9.10013 7.9998C8.82398 7.9998 8.60013 7.77595 8.60013 7.4998V6.0998C8.60013 5.82366 8.82398 5.5998 9.10013 5.5998Z" fill="currentColor"></path>
                      <path fill-rule="evenodd" clip-rule="evenodd" d="M4.20022 5.5998C4.47637 5.5998 4.70022 5.82366 4.70022 6.0998V7.4998C4.70022 7.77595 4.47637 7.9998 4.20022 7.9998C3.92408 7.9998 3.70022 7.77595 3.70022 7.4998V6.0998C3.70022 5.82366 3.92408 5.5998 4.20022 5.5998Z" fill="currentColor"></path>
                    </svg>
                  </div>
                  <span>Общение с GigaChat голосом и текстом</span>
                </li>
                <li class="login-page__feature">
                  <div aria-hidden="true" class="login-page__feature-icon">
                    <svg width="100%" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                      <path fill-rule="evenodd" clip-rule="evenodd" d="M5.19805 2.925C5.19805 1.86185 6.0599 1 7.12304 1H13.073C14.1362 1 14.998 1.86185 14.998 2.925V8.875C14.998 9.93815 14.1362 10.8 13.073 10.8H7.12305C6.88611 10.8 6.65917 10.7572 6.44954 10.6789L6.5286 10.8762H4.7576L5.30514 9.50972C5.23576 9.31099 5.19805 9.0974 5.19805 8.875V7.6499H5.1118L2.925 13.0399H3.9029L4.4111 11.7617H6.8828L7.3987 13.0399H8.4228L7.50132 10.8002H10.45V13.0749C10.45 14.138 9.58815 14.9999 8.525 14.9999H2.925C1.86185 14.9999 1 14.138 1 13.0749V7.4749C1 6.41175 1.86185 5.5499 2.925 5.5499H5.19805V2.925ZM7.12304 2.05C6.6398 2.05 6.24804 2.44175 6.24804 2.925V8.875C6.24804 9.35825 6.6398 9.75 7.12305 9.75H13.073C13.5563 9.75 13.948 9.35825 13.948 8.875V2.925C13.948 2.44175 13.5563 2.05 13.073 2.05H7.12304ZM10.0973 3.2751C10.3872 3.2751 10.6223 3.51015 10.6223 3.8001V8.0001C10.6223 8.29005 10.3872 8.5251 10.0973 8.5251C9.80731 8.5251 9.57226 8.29005 9.57226 8.0001V3.8001C9.57226 3.51015 9.80731 3.2751 10.0973 3.2751ZM12.1973 4.1501C12.4872 4.1501 12.7223 4.38515 12.7223 4.6751V7.1251C12.7223 7.41505 12.4872 7.6501 12.1973 7.6501C11.9073 7.6501 11.6723 7.41505 11.6723 7.1251V4.6751C11.6723 4.38515 11.9073 4.1501 12.1973 4.1501ZM8.52226 4.6751C8.52226 4.38515 8.28721 4.1501 7.99726 4.1501C7.70732 4.1501 7.47226 4.38515 7.47226 4.6751V7.1251C7.47226 7.41505 7.70732 7.6501 7.99726 7.6501C8.28721 7.6501 8.52226 7.41505 8.52226 7.1251V4.6751Z" fill="currentColor"></path>
                    </svg>
                  </div>
                  <span>Синтез и распознавание речи</span>
                </li>
                <li class="login-page__feature">
                  <div aria-hidden="true" class="login-page__feature-icon">
                    <svg width="100%" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                      <path fill-rule="evenodd" clip-rule="evenodd" d="M2.00183 8.00001C2.00183 4.68729 4.68732 2.0018 8.00004 2.0018C11.3128 2.0018 13.9982 4.68729 13.9982 8.00001C13.9982 11.3127 11.3128 13.9982 8.00004 13.9982C4.68732 13.9982 2.00183 11.3127 2.00183 8.00001ZM8.00004 1.0018C4.13504 1.0018 1.00183 4.13501 1.00183 8.00001C1.00183 11.865 4.13504 14.9982 8.00004 14.9982C11.865 14.9982 14.9982 11.865 14.9982 8.00001C14.9982 4.13501 11.865 1.0018 8.00004 1.0018ZM5.7159 7.29789C6.20089 7.29789 6.59404 6.90474 6.59404 6.41976C6.59404 5.93477 6.20089 5.54162 5.7159 5.54162C5.23092 5.54162 4.83777 5.93477 4.83777 6.41976C4.83777 6.90474 5.23092 7.29789 5.7159 7.29789ZM11.1604 6.41976C11.1604 6.90474 10.7673 7.29789 10.2823 7.29789C9.79733 7.29789 9.40417 6.90474 9.40417 6.41976C9.40417 5.93477 9.79733 5.54162 10.2823 5.54162C10.7673 5.54162 11.1604 5.93477 11.1604 6.41976ZM5.38981 9.50589C5.25174 9.26674 4.94594 9.1848 4.7068 9.32288C4.46765 9.46095 4.38571 9.76674 4.52379 10.0059C4.87596 10.6159 5.38249 11.1224 5.99248 11.4746C6.60246 11.8268 7.2944 12.0122 7.99875 12.0122C8.7031 12.0122 9.39504 11.8268 10.005 11.4746C10.615 11.1224 11.1215 10.6159 11.4737 10.0059C11.6118 9.76674 11.5299 9.46095 11.2907 9.32288C11.0516 9.1848 10.7458 9.26674 10.6077 9.50589C10.3433 9.96385 9.96299 10.3442 9.50502 10.6086C9.04706 10.873 8.52756 11.0122 7.99875 11.0122C7.46994 11.0122 6.95044 10.873 6.49248 10.6086C6.03451 10.3442 5.65422 9.96385 5.38981 9.50589Z" fill="currentColor"></path>
                    </svg>
                  </div>
                  <span>Создание персонализированных аватаров</span>
                </li>
              </ul>

              <#assign sberIdAlias = "sberid-oidc-rquid" />
              <#assign sberIdLoginUrl = "" />
              <#if social.providers??>
                <#list social.providers as p>
                  <#if p.alias == sberIdAlias || p.alias?contains("sberid")>
                    <#assign sberIdLoginUrl = p.loginUrl />
                  </#if>
                </#list>
              </#if>
              <#if !sberIdLoginUrl?has_content>
                <#assign sberIdLoginUrl = url.loginUrl + "?kc_idp_hint=" + sberIdAlias />
              </#if>

              <div class="login-page__actions">
                <a class="login-page__action login-page__action--primary" href="${sberIdLoginUrl}">Войти через Сбер ID</a>
                <button class="login-page__action" type="button" data-test-login-toggle="true" aria-controls="test-login-panel" aria-expanded="false">Тестовый вход</button>
                <button class="login-page__action" type="button">Войти по сертификату</button>
              </div>

              <#if realm.password>
                <div id="test-login-panel" class="test-login-panel" hidden>
                  <#if message?has_content && (message.type != 'warning' || !isAppInitiatedAction??)>
                    <div class="test-login-panel__error">
                      ${kcSanitize(message.summary)?no_esc}
                    </div>
                  </#if>

                  <form id="kc-form-login" class="test-login-form" action="${url.loginAction}" method="post">
                    <#if !usernameHidden??>
                      <label class="test-login-form__field" for="username">
                        <span>
                          <#if !realm.loginWithEmailAllowed>
                            Логин
                          <#elseif !realm.registrationEmailAsUsername>
                            Логин или почта
                          <#else>
                            Почта
                          </#if>
                        </span>
                        <input
                          tabindex="1"
                          id="username"
                          name="username"
                          value="${(login.username!'')}"
                          type="text"
                          autocomplete="username"
                          placeholder="Введите логин или email"
                          aria-invalid="<#if messagesPerField.existsError('username','password')>true</#if>"
                        />
                      </label>
                    </#if>

                    <label class="test-login-form__field" for="password">
                      <span>Пароль</span>
                      <input
                        tabindex="2"
                        id="password"
                        name="password"
                        type="password"
                        autocomplete="current-password"
                        placeholder="Введите пароль"
                        aria-invalid="<#if messagesPerField.existsError('username','password')>true</#if>"
                      />
                    </label>

                    <#if messagesPerField.existsError('username','password')>
                      <div class="test-login-panel__error">
                        ${kcSanitize(messagesPerField.getFirstError('username','password'))?no_esc}
                      </div>
                    </#if>

                    <div class="test-login-form__meta">
                      <#if realm.rememberMe && !usernameHidden??>
                        <label class="test-login-form__checkbox" for="rememberMe">
                          <input
                            tabindex="3"
                            id="rememberMe"
                            name="rememberMe"
                            type="checkbox"
                            <#if login.rememberMe??>checked</#if>
                          />
                          <span>${msg("rememberMe")}</span>
                        </label>
                      </#if>

                      <#if realm.resetPasswordAllowed>
                        <a class="test-login-form__link" tabindex="4" href="${url.loginResetCredentialsUrl}">
                          ${msg("doForgotPassword")}
                        </a>
                      </#if>
                    </div>

                    <input tabindex="5" class="test-login-form__submit" name="login" type="submit" value="Войти" />
                  </form>
                </div>
              </#if>

              <p class="login-page__agreement">
                Нажимая «Войти», вы принимаете<br />
                <a href="https://developers.sber.ru/docs/ru/policies/aistudio" target="_blank" rel="noopener noreferrer">Пользовательское соглашение</a>
              </p>
            </div>
          </main>

          <footer class="login-page__footer">
            <div class="login-page__footer-icon" aria-hidden="true">
              <svg width="36" viewBox="0 0 36 36" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path fill-rule="evenodd" clip-rule="evenodd" d="M18 20C15.6807 20 13.5025 21.107 11.9285 22.9959C10.5612 24.6366 9.73322 26.7614 9.54249 29H4.5C3.94772 29 3.5 29.4477 3.5 30C3.5 30.5523 3.94772 31 4.5 31H31.5C32.0523 31 32.5 30.5523 32.5 30C32.5 29.4477 32.0523 29 31.5 29H26.4575C26.2668 26.7614 25.4388 24.6366 24.0715 22.9959C22.4975 21.107 20.3193 20 18 20ZM24.4489 29C24.2641 27.2019 23.586 25.5373 22.5351 24.2762C21.2961 22.7894 19.659 22 18 22C16.3411 22 14.7039 22.7894 13.4649 24.2762C12.414 25.5373 11.7359 27.2019 11.5511 29H24.4489Z" fill="currentColor"></path>
                <path d="M10.1032 16.4709C9.72469 16.0687 9.09181 16.0495 8.68964 16.4281C8.28746 16.8066 8.26829 17.4394 8.6468 17.8416L10.1468 19.4354C10.5253 19.8375 11.1582 19.8567 11.5604 19.4782C11.9625 19.0997 11.9817 18.4668 11.6032 18.0646L10.1032 16.4709Z" fill="currentColor"></path>
                <path d="M3.5 25.2188C3.5 24.6665 3.94772 24.2188 4.5 24.2188H7.5C8.05228 24.2188 8.5 24.6665 8.5 25.2188C8.5 25.771 8.05228 26.2188 7.5 26.2188H4.5C3.94772 26.2188 3.5 25.771 3.5 25.2188Z" fill="currentColor"></path>
                <path d="M28.5 24.2188C27.9477 24.2188 27.5 24.6665 27.5 25.2188C27.5 25.771 27.9477 26.2188 28.5 26.2188H31.5C32.0523 26.2188 32.5 25.771 32.5 25.2188C32.5 24.6665 32.0523 24.2188 31.5 24.2188H28.5Z" fill="currentColor"></path>
                <path d="M27.3104 16.4281C27.7125 16.8066 27.7317 17.4394 27.3532 17.8416L25.8532 19.4354C25.4747 19.8375 24.8418 19.8567 24.4396 19.4782C24.0375 19.0997 24.0183 18.4668 24.3968 18.0646L25.8968 16.4709C26.2753 16.0687 26.9082 16.0495 27.3104 16.4281Z" fill="currentColor"></path>
                <path d="M19 13.0409V6C19 5.44772 18.5523 5 18 5C17.4477 5 17 5.44772 17 6V13.0409L14.2282 10.0959C13.8497 9.69371 13.2168 9.67453 12.8146 10.0531C12.4125 10.4316 12.3933 11.0644 12.7718 11.4666L17.2718 16.2479" fill="currentColor"></path>
                <path d="M17.2718 16.2479C17.4502 16.4132 17.7376 16.5625 18 16.5625C18.2883 16.5625 18.548 16.4405 18.7305 16.2454L23.2282 11.4666C23.6067 11.0644 23.5875 10.4316 23.1854 10.0531C22.7832 9.67453 22.1503 9.69371 21.7718 10.0959L19 13.0409" fill="currentColor"></path>
                <path d="M17.2741 16.2503C17.289 16.2661 17.2558 16.2331 17.2718 16.2479L17.2741 16.2503Z" fill="currentColor"></path>
              </svg>
            </div>
            <div class="login-page__privacy">
              Ознакомьтесь с <a href="https://www.sberbank.ru/privacy/policy#pdn" target="_blank" rel="noopener noreferrer">Политикой обработки персональных данных</a>
            </div>
          </footer>
        </div>
      </div>
    </div>

    <script>
      (function() {
        var toggle = document.querySelector('[data-test-login-toggle="true"]');
        var panel = document.getElementById('test-login-panel');
        if (!toggle || !panel) return;

        toggle.addEventListener('click', function() {
          var isHidden = panel.hasAttribute('hidden');
          if (isHidden) {
            panel.removeAttribute('hidden');
            toggle.setAttribute('aria-expanded', 'true');
            var firstInput = panel.querySelector('input[type="text"], input[type="email"]');
            if (firstInput) firstInput.focus();
          } else {
            panel.setAttribute('hidden', '');
            toggle.setAttribute('aria-expanded', 'false');
          }
        });

        if (panel.querySelector('.test-login-panel__error')) {
          panel.removeAttribute('hidden');
          toggle.setAttribute('aria-expanded', 'true');
        }
      })();
    </script>
  </#if>
</@layout.registrationLayout>
