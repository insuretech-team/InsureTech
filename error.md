forward-logs-shared.ts:95 Download the React DevTools for a better development experience: https://react.dev/link/react-devtools
forward-logs-shared.ts:95 [HMR] connected
button.tsx:52 A tree hydrated but some attributes of the server rendered HTML didn't match the client properties. This won't be patched up. This can happen if a SSR-ed Client Component used:

- A server/client branch `if (typeof window !== 'undefined')`.
- Variable input such as `Date.now()` or `Math.random()` which changes each time it's called.
- Date formatting in a user's locale which doesn't match the server.
- External changing data without sending a snapshot of it along with the HTML.
- Invalid HTML tag nesting.

It can also happen if the client has a browser extension installed which messes with the HTML before React loaded.

https://react.dev/link/hydration-mismatch

  ...
    <div className="ml-auto fl...">
      ...
        <MenuProvider scope={{Menu:[...], ...}} onClose={function Menu.useCallback} isUsingKeyboardRef={{current:false}} ...>
          <DropdownMenuTrigger asChild={true}>
            <DropdownMenuTrigger data-slot="dropdown-m..." asChild={true}>
              <MenuAnchor asChild={true} __scopeMenu={{Menu:[...], ...}}>
                <PopperAnchor __scopePopper={{Menu:[...], ...}} asChild={true} ref={null}>
                  <Primitive.div asChild={true} ref={function}>
                    <Primitive.div.Slot ref={function}>
                      <Primitive.div.SlotClone ref={function}>
                        <Primitive.button type="button" id="radix-_R_c..." aria-haspopup="menu" aria-expanded={false} ...>
                          <Primitive.button.Slot type="button" id="radix-_R_c..." aria-haspopup="menu" ...>
                            <Primitive.button.SlotClone type="button" id="radix-_R_c..." aria-haspopup="menu" ...>
                              <Button variant="outline" size="icon" className="relative s..." type="button" ...>
                                <button
                                  data-slot="dropdown-menu-trigger"
                                  data-variant="outline"
                                  data-size="icon"
                                  className={"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded..."}
                                  type="button"
+                                 id="radix-_R_ctbn5rlb_"
-                                 id="radix-_R_1jpbn5rlb_"
                                  aria-haspopup="menu"
                                  aria-expanded={false}
                                  aria-controls={undefined}
                                  data-state="closed"
                                  data-disabled={undefined}
                                  disabled={false}
                                  onPointerDown={function handleEvent}
                                  onKeyDown={function handleEvent}
                                  ref={function}
                                >
          ...
      ...
        <MenuProvider scope={{Menu:[...], ...}} onClose={function Menu.useCallback} isUsingKeyboardRef={{current:false}} ...>
          <DropdownMenuTrigger asChild={true}>
            <DropdownMenuTrigger data-slot="dropdown-m..." asChild={true}>
              <MenuAnchor asChild={true} __scopeMenu={{Menu:[...], ...}}>
                <PopperAnchor __scopePopper={{Menu:[...], ...}} asChild={true} ref={null}>
                  <Primitive.div asChild={true} ref={function}>
                    <Primitive.div.Slot ref={function}>
                      <Primitive.div.SlotClone ref={function}>
                        <Primitive.button type="button" id="radix-_R_k..." aria-haspopup="menu" aria-expanded={false} ...>
                          <Primitive.button.Slot type="button" id="radix-_R_k..." aria-haspopup="menu" ...>
                            <Primitive.button.SlotClone type="button" id="radix-_R_k..." aria-haspopup="menu" ...>
                              <Button variant="link" size="icon" className="size-9 rou..." type="button" ...>
                                <button
                                  data-slot="dropdown-menu-trigger"
                                  data-variant="link"
                                  data-size="icon"
                                  className={"inline-flex items-center justify-center gap-2 whitespace-nowrap text-sm..."}
                                  type="button"
+                                 id="radix-_R_ktbn5rlb_"
-                                 id="radix-_R_2jpbn5rlb_"
                                  aria-haspopup="menu"
                                  aria-expanded={false}
                                  aria-controls={undefined}
                                  data-state="closed"
                                  data-disabled={undefined}
                                  disabled={false}
                                  onPointerDown={function handleEvent}
                                  onKeyDown={function handleEvent}
                                  ref={function}
                                >
          ...

error @ intercept-console-error.ts:42
(anonymous) @ react-dom-client.development.js:5731
runWithFiberInDEV @ react-dom-client.development.js:986
emitPendingHydrationWarnings @ react-dom-client.development.js:5730
completeWork @ react-dom-client.development.js:12862
runWithFiberInDEV @ react-dom-client.development.js:989
completeUnitOfWork @ react-dom-client.development.js:19133
performUnitOfWork @ react-dom-client.development.js:19014
workLoopConcurrentByScheduler @ react-dom-client.development.js:18991
renderRootConcurrent @ react-dom-client.development.js:18973
performWorkOnRoot @ react-dom-client.development.js:17834
performWorkOnRootViaSchedulerTask @ react-dom-client.development.js:20384
performWorkUntilDeadline @ scheduler.development.js:45
<button>
exports.jsxDEV @ react-jsx-dev-runtime.development.js:342
Button @ button.tsx:52
react_stack_bottom_frame @ react-dom-client.development.js:28038
renderWithHooksAgain @ react-dom-client.development.js:8084
renderWithHooks @ react-dom-client.development.js:7996
updateFunctionComponent @ react-dom-client.development.js:10501
beginWork @ react-dom-client.development.js:12136
runWithFiberInDEV @ react-dom-client.development.js:986
performUnitOfWork @ react-dom-client.development.js:18997
workLoopConcurrentByScheduler @ react-dom-client.development.js:18991
renderRootConcurrent @ react-dom-client.development.js:18973
performWorkOnRoot @ react-dom-client.development.js:17834
performWorkOnRootViaSchedulerTask @ react-dom-client.development.js:20384
performWorkUntilDeadline @ scheduler.development.js:45
<Button>
exports.jsxDEV @ react-jsx-dev-runtime.development.js:342
DashboardHeader @ dashboard-header.tsx:74
react_stack_bottom_frame @ react-dom-client.development.js:28038
renderWithHooksAgain @ react-dom-client.development.js:8084
renderWithHooks @ react-dom-client.development.js:7996
updateFunctionComponent @ react-dom-client.development.js:10501
beginWork @ react-dom-client.development.js:12136
runWithFiberInDEV @ react-dom-client.development.js:986
performUnitOfWork @ react-dom-client.development.js:18997
workLoopConcurrentByScheduler @ react-dom-client.development.js:18991
renderRootConcurrent @ react-dom-client.development.js:18973
performWorkOnRoot @ react-dom-client.development.js:17834
performWorkOnRootViaSchedulerTask @ react-dom-client.development.js:20384
performWorkUntilDeadline @ scheduler.development.js:45
<DashboardHeader>
exports.jsxDEV @ react-jsx-dev-runtime.development.js:342
DashboardLayout @ dashboard-layout.tsx:39
react_stack_bottom_frame @ react-dom-client.development.js:28038
renderWithHooksAgain @ react-dom-client.development.js:8084
renderWithHooks @ react-dom-client.development.js:7996
updateFunctionComponent @ react-dom-client.development.js:10501
beginWork @ react-dom-client.development.js:12136
runWithFiberInDEV @ react-dom-client.development.js:986
performUnitOfWork @ react-dom-client.development.js:18997
workLoopConcurrentByScheduler @ react-dom-client.development.js:18991
renderRootConcurrent @ react-dom-client.development.js:18973
performWorkOnRoot @ react-dom-client.development.js:17834
performWorkOnRootViaSchedulerTask @ react-dom-client.development.js:20384
performWorkUntilDeadline @ scheduler.development.js:45
<DashboardLayout>
exports.jsxDEV @ react-jsx-dev-runtime.development.js:342
EmployeesPage @ employees-table.tsx:58
react_stack_bottom_frame @ react-dom-client.development.js:28038
renderWithHooksAgain @ react-dom-client.development.js:8084
renderWithHooks @ react-dom-client.development.js:7996
updateFunctionComponent @ react-dom-client.development.js:10501
beginWork @ react-dom-client.development.js:12085
replayBeginWork @ react-dom-client.development.js:19054
runWithFiberInDEV @ react-dom-client.development.js:986
replaySuspendedUnitOfWork @ react-dom-client.development.js:19018
renderRootConcurrent @ react-dom-client.development.js:18901
performWorkOnRoot @ react-dom-client.development.js:17834
performWorkOnRootViaSchedulerTask @ react-dom-client.development.js:20384
performWorkUntilDeadline @ scheduler.development.js:45
"use client"
page @ page.tsx:4
initializeElement @ react-server-dom-turbopack-client.browser.development.js:1941
(anonymous) @ react-server-dom-turbopack-client.browser.development.js:4623
initializeModelChunk @ react-server-dom-turbopack-client.browser.development.js:1828
getOutlinedModel @ react-server-dom-turbopack-client.browser.development.js:2337
parseModelString @ react-server-dom-turbopack-client.browser.development.js:2729
(anonymous) @ react-server-dom-turbopack-client.browser.development.js:4554
initializeModelChunk @ react-server-dom-turbopack-client.browser.development.js:1828
resolveModelChunk @ react-server-dom-turbopack-client.browser.development.js:1672
processFullStringRow @ react-server-dom-turbopack-client.browser.development.js:4442
processFullBinaryRow @ react-server-dom-turbopack-client.browser.development.js:4300
processBinaryChunk @ react-server-dom-turbopack-client.browser.development.js:4523
progress @ react-server-dom-turbopack-client.browser.development.js:4799
<page>
Promise.all @ VM3983 <anonymous>:1
Promise.all @ VM3983 <anonymous>:1
initializeFakeTask @ react-server-dom-turbopack-client.browser.development.js:3390
initializeDebugInfo @ react-server-dom-turbopack-client.browser.development.js:3415
initializeDebugChunk @ react-server-dom-turbopack-client.browser.development.js:1772
processFullStringRow @ react-server-dom-turbopack-client.browser.development.js:4389
processFullBinaryRow @ react-server-dom-turbopack-client.browser.development.js:4300
processBinaryChunk @ react-server-dom-turbopack-client.browser.development.js:4523
progress @ react-server-dom-turbopack-client.browser.development.js:4799
"use server"
ResponseInstance @ react-server-dom-turbopack-client.browser.development.js:2784
createResponseFromOptions @ react-server-dom-turbopack-client.browser.development.js:4660
exports.createFromReadableStream @ react-server-dom-turbopack-client.browser.development.js:5064
module evaluation @ app-index.tsx:211
(anonymous) @ dev-base.ts:244
runModuleExecutionHooks @ dev-base.ts:278
instantiateModule @ dev-base.ts:238
getOrInstantiateModuleFromParent @ dev-base.ts:162
commonJsRequire @ runtime-utils.ts:389
(anonymous) @ app-next-turbopack.ts:11
(anonymous) @ app-bootstrap.ts:79
loadScriptsInSequence @ app-bootstrap.ts:23
appBootstrap @ app-bootstrap.ts:61
module evaluation @ app-next-turbopack.ts:10
(anonymous) @ dev-base.ts:244
runModuleExecutionHooks @ dev-base.ts:278
instantiateModule @ dev-base.ts:238
getOrInstantiateRuntimeModule @ dev-base.ts:128
registerChunk @ runtime-backend-dom.ts:57
await in registerChunk
registerChunk @ dev-base.ts:1149
(anonymous) @ dev-backend-dom.ts:126
(anonymous) @ dev-backend-dom.ts:126
employees-table.tsx:27  GET http://localhost:3000/api/employees 404 (Not Found)
loadEmployees @ employees-table.tsx:27
EmployeesPage.useEffect @ employees-table.tsx:51
react_stack_bottom_frame @ react-dom-client.development.js:28123
runWithFiberInDEV @ react-dom-client.development.js:986
commitHookEffectListMount @ react-dom-client.development.js:13692
commitHookPassiveMountEffects @ react-dom-client.development.js:13779
commitPassiveMountOnFiber @ react-dom-client.development.js:16733
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16753
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16753
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16753
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16753
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16753
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16753
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16753
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16753
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16753
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16725
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:17010
recursivelyTraversePassiveMountEffects @ react-dom-client.development.js:16678
commitPassiveMountOnFiber @ react-dom-client.development.js:16768
flushPassiveEffects @ react-dom-client.development.js:19859
flushPendingEffects @ react-dom-client.development.js:19785
performSyncWorkOnRoot @ react-dom-client.development.js:20396
flushSyncWorkAcrossRoots_impl @ react-dom-client.development.js:20241
flushSpawnedWork @ react-dom-client.development.js:19752
commitRoot @ react-dom-client.development.js:19335
commitRootWhenReady @ react-dom-client.development.js:18178
performWorkOnRoot @ react-dom-client.development.js:18054
performWorkOnRootViaSchedulerTask @ react-dom-client.development.js:20384
performWorkUntilDeadline @ scheduler.development.js:45
"use client"
page @ page.tsx:4
initializeElement @ react-server-dom-turbopack-client.browser.development.js:1941
(anonymous) @ react-server-dom-turbopack-client.browser.development.js:4623
initializeModelChunk @ react-server-dom-turbopack-client.browser.development.js:1828
getOutlinedModel @ react-server-dom-turbopack-client.browser.development.js:2337
parseModelString @ react-server-dom-turbopack-client.browser.development.js:2729
(anonymous) @ react-server-dom-turbopack-client.browser.development.js:4554
initializeModelChunk @ react-server-dom-turbopack-client.browser.development.js:1828
resolveModelChunk @ react-server-dom-turbopack-client.browser.development.js:1672
processFullStringRow @ react-server-dom-turbopack-client.browser.development.js:4442
processFullBinaryRow @ react-server-dom-turbopack-client.browser.development.js:4300
processBinaryChunk @ react-server-dom-turbopack-client.browser.development.js:4523
progress @ react-server-dom-turbopack-client.browser.development.js:4799
<page>
Promise.all @ VM3983 <anonymous>:1
Promise.all @ VM3983 <anonymous>:1
initializeFakeTask @ react-server-dom-turbopack-client.browser.development.js:3390
initializeDebugInfo @ react-server-dom-turbopack-client.browser.development.js:3415
initializeDebugChunk @ react-server-dom-turbopack-client.browser.development.js:1772
processFullStringRow @ react-server-dom-turbopack-client.browser.development.js:4389
processFullBinaryRow @ react-server-dom-turbopack-client.browser.development.js:4300
processBinaryChunk @ react-server-dom-turbopack-client.browser.development.js:4523
progress @ react-server-dom-turbopack-client.browser.development.js:4799
"use server"
ResponseInstance @ react-server-dom-turbopack-client.browser.development.js:2784
createResponseFromOptions @ react-server-dom-turbopack-client.browser.development.js:4660
exports.createFromReadableStream @ react-server-dom-turbopack-client.browser.development.js:5064
module evaluation @ app-index.tsx:211
(anonymous) @ dev-base.ts:244
runModuleExecutionHooks @ dev-base.ts:278
instantiateModule @ dev-base.ts:238
getOrInstantiateModuleFromParent @ dev-base.ts:162
commonJsRequire @ runtime-utils.ts:389
(anonymous) @ app-next-turbopack.ts:11
(anonymous) @ app-bootstrap.ts:79
loadScriptsInSequence @ app-bootstrap.ts:23
appBootstrap @ app-bootstrap.ts:61
module evaluation @ app-next-turbopack.ts:10
(anonymous) @ dev-base.ts:244
runModuleExecutionHooks @ dev-base.ts:278
instantiateModule @ dev-base.ts:238
getOrInstantiateRuntimeModule @ dev-base.ts:128
registerChunk @ runtime-backend-dom.ts:57
await in registerChunk
registerChunk @ dev-base.ts:1149
(anonymous) @ dev-backend-dom.ts:126
(anonymous) @ dev-backend-dom.ts:126
forward-logs-shared.ts:95 [Fast Refresh] rebuilding
forward-logs-shared.ts:95 [Fast Refresh] done in 1267ms
forward-logs-shared.ts:95 [Fast Refresh] rebuilding
forward-logs-shared.ts:95 [Fast Refresh] done in 1514ms
forward-logs-shared.ts:95 [Fast Refresh] rebuilding
forward-logs-shared.ts:95 [Fast Refresh] done in 1579ms
forward-logs-shared.ts:95 [Fast Refresh] rebuilding
forward-logs-shared.ts:95 [Fast Refresh] done in 781ms
forward-logs-shared.ts:95 [Fast Refresh] rebuilding
forward-logs-shared.ts:95 [Fast Refresh] done in 1207ms
forward-logs-shared.ts:95 [Fast Refresh] rebuilding
forward-logs-shared.ts:95 [Fast Refresh] done in 6089ms
forward-logs-shared.ts:95 Image with src "/logos/insuretech-brand.png" was detected as the Largest Contentful Paint (LCP). Please add the `loading="eager"` property if this image is above the fold.
Read more: https://nextjs.org/docs/app/api-reference/components/image#loading
warn @ forward-logs-shared.ts:95
warnOnce @ warn-once.ts:6
(anonymous) @ get-img-props.ts:647



