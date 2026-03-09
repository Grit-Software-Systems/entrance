## Purpose
A Go package for authenticating Entra External ID users using native authentication, https://learn.microsoft.com/en-us/entra/identity-platform/reference-native-authentication-api

## Coding Guidelines

### Core Principles (Inviolate)

- Never nest function calls. They make code harder to debug. To see the output of the nested function, one has to step into the nesting function and observe the output from the nesting function that was passed to the nesting function as an argument.
- Never return a function call. Rather assign the output of the function call to a variable named `result` and return `result`. This eases debugging: one can set a breakpoint on the return statement to see the value of the `result` variable.
- Name the variable returned by a function `result`. Then readers of the code will readily be able to see which variable will ultimately be returned by a function.
- Always declare the types of variables in any language that permits it. Therefore, in Python, always use type hints, even though those are optional. In a language like C# that is generally typed but which permits duck typing using `var`, one should never, ever use the duck typing option because then the reader has to mentally compute the implied type.
- Never, ever use abbreviations or acronyms in names. Verbose naming helps to make code self-documenting.
- Use comments only to explain code that is otherwise not self-explanatory. Of course, comments are required for tools to be used by agents; that's an exception.
- Try very hard to limit the length of any subroutine to 40 lines.
- Generally, do not use literals in code. Rather declare a constant, the name of which explains the meaning of the literal. The exceptions are Boolean literals and the integers used for indexing and index arithmetic.
- Group the declarations of the same types of members/things together: constants must be grouped together, functions must be grouped separately from properties and so on.
- Arrange things in the following order:
    - Private constants in alphabetical order
    - Public constants in alphabetical order
    - Private static variables in alphabetical order
    - Public static variables in alphabetical order
    - Constructors
    - Private interface definitions in alphabetical order
    - Properties, public and private in alphabetical order
    - Functions, public and private in alphabetical order
    - Private classes, structs and enums in alphabetical order
- Within a group, arrange strictly alphabetically; for example, group functions together and arrange them alphabetically, and group properties together, and arrange them alphabetically. In Python, this convention has the handy side effect of making all private things, which would be prefixed with an underscore, come before all public things.
- There should only be one public class in a file.
- Import/using statements of built-in libraries should precede imports of third-party and local libraries. Within those groups, the statements should be arranged alphabetically.
- LINQ statements should be written in LINQ Query syntax rather than LINQ method syntax.
- In Python comprehensions, LINQ query statements, and similar constructs in other languages, the iterated variable should be named `item`.
- High test coverage should be achieved without ever making a class/member/property/function in product code visible (public) solely for testing. The visibility of anything should be determined by the design of the product, defaulting to the lowest visibility necessary. For testing, one can use reflection, add product code files to test projects as linked files, make test projects friends of product classes, and so on.

---
### Web Service Coding Guidelines
Web services consist of business logic wrapped by a framework, such as FastAPI or Express or ASP.NET.
Web services should always be programmed to minimize the code that is any file that depends on the framework.
Code that has any references to the framework should simply authorize a request and a route it to the business layer which should be in an independently-testable library/package.
---
### Azure Functions
Any Azure Functions should be designed so to minimize the code that is any file that depends on Azure Function libraries. That code should quickly route processing to a separate, independently testable library/package that does the actual work.
---

### JavaScript Coding Guidelines

#### Language and Runtime

- Use modern ECMAScript syntax supported by the target runtime.
- Use `const` by default. Use `let` only when reassignment is required. Never use `var`.
- Enable strict mode where applicable.
- Do not use TypeScript.

#### Typing and Structure

- When type systems are not available, simulate type clarity through:
    - Explicit runtime validation.
    - Well-named variables and parameters.
    - JSDoc comments when public APIs require clarity.
- Avoid dynamic shape mutation of objects after construction.
- Prefer explicit object construction over ad hoc property assignment.

#### Functions

- Follow the core rule: never nest function calls and never return a function call.
- Avoid anonymous inline callbacks except when structurally unavoidable (for example, event listeners).
- Prefer named functions over anonymous functions for improved stack traces.
- Limit functions to 40 lines whenever possible.
- Avoid deeply nested control flow. Extract intermediate results into clearly named variables.

#### Error Handling

- Always handle rejected Promises.
- Use `try`/`catch` around `await` statements when failure is possible.
- Never swallow exceptions silently.
- Error objects must contain meaningful messages and contextual data.

#### Asynchronous Code

- Prefer `async`/`await` over chained Promise syntax.
- Avoid chaining multiple `.then()` calls.
- Assign awaited results to explicitly typed variables before further use.
- Do not mix callback patterns with Promise-based patterns.

#### Modules and Imports

- Built-in modules first, third-party modules second, local modules third.
- Alphabetize imports within each group.
- Avoid wildcard imports.
- Export only what is necessary.

#### Object-Oriented and Functional Structure

- Prefer classes when modeling stateful domain concepts.
- Avoid overly functional, point-free styles that obscure intermediate values.
- Avoid mutation of shared state; if mutation is required, isolate it clearly.

#### Constants

- Replace string literals used more than once with named constants.
- Use uppercase for constants.
- Avoid magic numbers.

#### State Management

- State management does not rely on React.useState, but rather on a state reducer, which should be implemented in the src\lib\state.js module.
- App.jsx would have the following:
```
  import { useTranslation } from "react-i18next";
  import {
      initializeState, reduceState
  } from "@/lib/state.js";
  import {
    RESOURCE_FILE_DEFAULT
} from "@/lib/constants.js";

  export const App = () => {
    const { t:translate } = useTranslation(RESOURCE_FILE_DEFAULT);
    const { i18n:translations } = useTranslation();
    const [state, dispatchState] = useReducer(reduceState, undefined, initializeState);
    const reference = useRef({cancelled: false});
```
- Then translate, translations, state, dispatchState and reference are passed to all other components:
```
   <MyOtherComponent  state={state}
                      dispatchState={dispatchState}
                      reference={reference}
                      translate={translate}
                      translations={translations}  
    />
```
- When a component calls dispatchState, it does so like this:
```
    const onInitiateImport = () => dispatchState({
          type: ACTION_ON_SOMETHING_HAPPENED,
          payload: myValue
      });
```
    Here, ACTION_ON_SOMETHING_HAPPENED is a constant that would be declared in src\lib\state.js

#### Globalization

Components use the translate and translations functions for globalization. Any natural language string displayed to the user is displayed like so:
'''
{translate("LABEL_IMPORT_BETS")}
'''
Here, LABEL_IMPORT_BETS is a constant in the src\locales\common.json file.

#### Logging

All logging is done via the log(), report() and warn() functions exported from src\lib\utilities.js.
- Call log() to log information.
- Call warn() to log warnings.
- Call report() to log errors.

#### Testing

- Test through public surface area whenever possible.
- Do not alter visibility solely for testing.
- Use dependency injection rather than monkey patching.

---

### C# Coding Guidelines

#### Language Usage

- Never use `var`. Always explicitly declare types.
- Use nullable reference types and enable nullable context.
- Avoid dynamic typing (`dynamic`) unless absolutely required by interop constraints.
- Prefer explicit access modifiers; never rely on defaults.

#### Naming

- Use PascalCase for types, methods, and properties.
- Use camelCase for local variables and parameters.
- Do not use abbreviations or acronyms.
- Avoid suffixes such as `Mgr`, `Svc`, `Util`.

#### Classes and Structure

- One public class per file.
- Keep classes cohesive and focused on a single responsibility.
- Prefer composition over inheritance unless inheritance models a true "is-a" relationship.
- Avoid partial classes unless required by tooling.

#### Methods

- Do not nest method calls.
- Do not return method calls.
- Use the `result` variable for return values.
- Avoid long parameter lists; prefer parameter objects when needed.
- Limit methods to 40 lines.

#### LINQ

- Always use LINQ query syntax.
- Name iterated variables `item`.
- Extract intermediate query results into variables before projection or materialization.

#### Exceptions

- Throw specific exception types.
- Do not catch `Exception` unless rethrowing with additional context.
- Use `throw;` to preserve stack traces.
- Avoid using exceptions for control flow.

#### Asynchronous Programming

- Use `async`/`await`.
- Avoid `.Result` or `.Wait()`.
- Return `Task` or `Task<T>` explicitly.
- Avoid `async void` except for event handlers.

#### Dependency Injection

- Use constructor injection.
- Avoid service locator patterns.
- Keep constructors simple and focused on dependency assignment.

#### Constants and Configuration

- Replace repeated literals with `const` or `static readonly` fields.
- Avoid embedding configuration values in code.
- Use strongly typed configuration binding.

#### Testing

- Use internal visibility with friend assemblies if required.
- Avoid making members public solely for testing.
- Test behavior, not implementation details.

---

### .NET Architectural Guidelines

#### Project Structure

- Separate concerns into distinct projects when appropriate:
    - Domain
    - Application
    - Infrastructure
    - Presentation
- Avoid circular dependencies.
- Domain layer must not depend on infrastructure.

#### Dependency Management

- Prefer built-in dependency injection container.
- Register services explicitly.
- Avoid global static state.

#### Configuration

- Use strongly typed options patterns.
- Validate configuration at startup.
- Avoid accessing configuration directly in domain logic.

#### Logging

- Use structured logging.
- Avoid string concatenation in log messages.
- Inject logging abstractions rather than referencing static loggers.

#### Data Access

- Use repositories only when abstraction adds value.
- Keep data access logic isolated.
- Avoid leaking persistence concerns into domain models.

#### API Design

- Return explicit types.
- Avoid returning anonymous objects.
- Validate inputs at boundaries.
- Use meaningful HTTP status codes.

#### Security

- Never store secrets in source code.
- Use secure configuration providers.
- Validate external inputs rigorously.

#### Performance

- Avoid premature optimization.
- Measure before optimizing.
- Use asynchronous I/O where appropriate.

### Python Coding Guidelines

#### Language and Runtime

- Support only the Python versions explicitly targeted by the product.
- Run with warnings enabled in CI when feasible.
- Prefer the standard library over third-party libraries unless the third-party library is clearly superior and stable.

#### Typing and Interfaces

- Always use type hints for:
    - Function parameters and return values.
    - Local variables, including loop variables when their type is not obvious.
    - Class members, including instance attributes and class variables.
- Prefer `typing` and `collections.abc` types (`Iterable`, `Mapping`, `Sequence`) over concrete container types when callers should not depend on the implementation.
- Prefer `@dataclass` for simple data carriers, with explicit type annotations on all fields.
- Avoid `Any`. When unavoidable (interop, dynamic payloads), constrain it at the boundary and convert to typed structures immediately.
- Validate untrusted inputs at module boundaries and convert them into typed domain objects.

#### Naming

- Never use abbreviations or acronyms. Use verbose, descriptive names.
- Follow PEP 8 casing conventions while remaining verbose:
    - Modules: `lowercase_with_underscores`
    - Classes: `PascalCase`
    - Functions and variables: `lowercase_with_underscores`
    - Constants: `UPPERCASE_WITH_UNDERSCORES`
- Avoid single-letter names except for indexing variables (and only where indexing is actually occurring).

#### Files and Imports

- One public class per file (consistent with the core guideline).
- Import order:
    - Python standard library imports first.
    - Third-party imports second.
    - Local application imports third.
- Alphabetize imports within each group.
- Avoid `from module import *`.
- Avoid importing heavy dependencies at module import time unless necessary; prefer lazy imports inside functions when they reduce startup time and do not violate the “no nested calls” rule.

#### Functions and Control Flow

- Never nest function calls.
- Never return a function call; assign to `result` and return `result`.
- Name the variable returned by a function `result`.
- Keep functions under 40 lines whenever possible by extracting cohesive helpers.
- Prefer early returns and guard clauses to reduce indentation and improve readability.
- Avoid deep nesting of `if`/`for`/`while` blocks; refactor into smaller functions.

#### Comprehensions and Iteration

- In comprehensions and similar constructs, the iterated variable must be named `item`.
- Prefer clarity over cleverness:
    - Use list/dict/set comprehensions only when the expression is simple and readable.
    - If the comprehension becomes complex, expand into explicit loops with intermediate variables.
- Avoid generator expressions as arguments to another call (because that creates nested evaluation). Prefer assigning the generator to a variable first.

#### Exceptions and Error Handling

- Raise specific exception types.
- Never use bare `except:`. Catch explicit exception types.
- Do not swallow exceptions. If catching for context, rethrow with additional context and preserve the original exception (`raise ... from exception`).
- Use exceptions for exceptional conditions, not normal control flow.

#### Resource Management

- Use context managers (`with`) for resources (files, sockets, locks, database connections).
- Prefer `contextlib` utilities (`@contextmanager`, `ExitStack`) for composing resource lifetimes.
- Ensure resources are closed in `finally` blocks when a context manager is not available.

#### Asynchronous Python

- Prefer `async`/`await` for I/O bound concurrency when the ecosystem supports it.
- Do not mix blocking I/O in the event loop.
- Always `await` tasks you create, or explicitly manage them (and handle exceptions).
- Keep async call stacks readable by assigning intermediate awaited results to named variables.

#### Data Modeling

- Prefer small, explicit domain objects over loose dictionaries.
- When dictionaries are required (JSON payloads), convert them to typed objects as early as possible.
- Avoid mutating shared structures across layers; return new values where practical.

#### Constants and Literals

- Avoid literals. Use named constants that explain meaning.
- Exceptions remain:
    - Boolean literals.
    - Integers used for indexing and index arithmetic.
- Prefer module-level constants for shared values; keep them grouped and alphabetized per the core ordering rules.

#### Logging

- Use a consistent logging framework (standard `logging` unless otherwise specified).
- Prefer structured, contextual logging via:
    - Explicit fields in log messages.
    - Extra context dictionaries where supported.
- Avoid string concatenation in log messages; use parameterized logging.
- Ensure logs avoid secrets and sensitive data.

#### Testing

- High coverage without changing product visibility solely for tests.
- Prefer testing public behavior and module boundaries.
- Use:
    - Dependency injection via parameters or constructor injection.
    - Fakes/mocks at boundaries.
    - Monkeypatching only as a last resort and only in tests.

#### Formatting and Style

- Use an auto-formatter (for example, Black) consistently across the repo if adopted by the team.
- Keep line lengths consistent with the chosen formatter policy.
- Comments are only for non-obvious reasoning, tradeoffs, and constraints—not restating code.

#### Dependency and Packaging Hygiene

- Pin dependencies with explicit versions for repeatable builds.
- Keep runtime dependencies minimal.
- Separate development dependencies from production dependencies.
- Validate license compatibility for third-party libraries.

#### Summary

Python code must prioritize **debuggability**, **explicitness**, and **traceability**:
- No nested calls.
- No returning calls.
- Always `result`.
- Types everywhere.
- Readable, flat control flow.


---

### Closing Principle

These platform and language-specific rules supplement, but never override, the core guidelines. Debuggability, explicitness, and structural clarity take precedence over brevity or stylistic trends.

